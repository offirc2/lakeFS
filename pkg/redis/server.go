package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/treeverse/lakefs/pkg/logging"

	"github.com/danwakefield/fnmatch"
	"github.com/treeverse/lakefs/pkg/graveler"
	"github.com/treeverse/lakefs/pkg/ident"
	"github.com/treeverse/lakefs/pkg/resp"
)

const (
	redisClientNameContextKey = "redis_client_name_ctx_key"
)

type GravelerRedisAPI interface {
	graveler.VersionController
	graveler.KeyValueStore
}

func extractVersion(store GravelerRedisAPI, ctx context.Context) (record *graveler.RepositoryRecord, ref graveler.Ref, err error) {
	val := ctx.Value(redisClientNameContextKey)
	if val == nil {
		return nil, "",
			fmt.Errorf("error: no repository and ref in client name. Call CLIENT SETNAME")
	}
	clientName := val.(string)
	atPartIndex := strings.Index(clientName, "@")
	if atPartIndex == -1 {
		atPartIndex = 0
	}
	clientName = clientName[atPartIndex+1:]
	parts := strings.SplitN(clientName, "/", 2)
	repoId := parts[0]
	ref = graveler.Ref(parts[1])
	repo, err := store.GetRepository(ctx, graveler.RepositoryID(repoId))
	if err != nil {
		return nil, "", err
	}
	if _, err := store.Dereference(ctx, repo, ref); err != nil {
		return nil, "", fmt.Errorf("reference not found: %s", ref)
	}
	return repo, ref, nil
}

func Commands(store GravelerRedisAPI, log logging.Logger) resp.Handler {
	router := resp.NewRouter()

	router.Add(&resp.Command{
		Name:     "client setname",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			ctx := context.WithValue(request.Context(), redisClientNameContextKey, string(args[0]))
			request.SetContext(ctx)
			w.WriteOK()
		},
	})

	router.Add(&resp.Command{
		Name:     "client getname",
		Validate: resp.NoArgs(),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			ctx := request.Context()
			val := ctx.Value(redisClientNameContextKey)
			if val == nil {
				w.WriteNull()
				return
			}
			clientName := val.(string)
			w.WriteSimpleString(clientName)
		},
	})
	router.Add(&resp.Command{
		Name:     "get",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			value, err := store.Get(request.Context(), repo, ref, args[0])
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			w.WriteBulkString(value.Data)
		},
	})
	router.Add(&resp.Command{
		Name:     "keys",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			pattern := string(args[0])
			iter, err := store.List(request.Context(), repo, ref, 1024)
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			keys := make([]string, 0)
			for iter.Next() {
				key := iter.Value().Key.String()
				if fnmatch.Match(pattern, key, fnmatch.FNM_FILE_NAME) {
					keys = append(keys, key)
				}
			}
			w.WriteArray(len(keys))
			for _, k := range keys {
				w.WriteSimpleString(k)
			}
		},
	})
	router.Add(&resp.Command{
		Name:     "set",
		Validate: resp.NArgs(2),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			err = store.Set(request.Context(), repo, graveler.BranchID(ref), args[0], graveler.Value{
				Identity: ident.NewAddressWriter().MarshalBytes(args[1]).Identity(),
				Data:     args[1],
			})
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			w.WriteOK()
		},
	})
	router.Add(&resp.Command{
		Name:     "del",
		Validate: resp.MinArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			keys := make([]graveler.Key, 0)
			for _, k := range args {
				keys = append(keys, k)
			}
			err = store.DeleteBatch(request.Context(), repo, graveler.BranchID(ref), keys)
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			w.WriteInteger(len(keys))
		},
	})
	router.Add(&resp.Command{
		Name:     "checkout",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			_, _, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			branch := string(args[0])
			// checkout
			clientName := request.Context().Value(redisClientNameContextKey).(string)
			refStartsAt := strings.LastIndex(clientName, "/") + 1
			clientName = clientName[0:refStartsAt] + branch
			request.SetContext(context.WithValue(request.Context(), redisClientNameContextKey, clientName))
			w.WriteOK()
		},
	})
	router.Add(&resp.Command{
		Name:     "branch create",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			branch := string(args[0])
			_, err = store.CreateBranch(request.Context(), repo, graveler.BranchID(branch), ref)
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			// checkout
			clientName := request.Context().Value(redisClientNameContextKey).(string)
			refStartsAt := strings.LastIndex(clientName, "/") + 1
			clientName = clientName[0:refStartsAt] + branch
			request.SetContext(context.WithValue(request.Context(), redisClientNameContextKey, clientName))
			w.WriteOK()
		},
	})
	router.Add(&resp.Command{
		Name:     "branch reset",
		Validate: resp.NoArgs(),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			err = store.Reset(request.Context(), repo, graveler.BranchID(ref))
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			w.WriteOK()
		},
	})
	router.Add(&resp.Command{
		Name:     "diff",
		Validate: resp.ArgsBetween(0, 1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			var diffIter graveler.DiffIterator
			if len(args) == 1 {
				// diff with another branch
				diffIter, err = store.Diff(request.Context(), repo, ref, graveler.Ref(args[0]))
			} else {
				// diff current uncommitted
				diffIter, err = store.DiffUncommitted(request.Context(), repo, graveler.BranchID(ref))
			}
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			results := make([][2]string, 0)
			for diffIter.Next() {
				current := diffIter.Value().Copy()
				diffType := ""
				switch current.Type {
				case graveler.DiffTypeAdded:
					diffType = "added"
				case graveler.DiffTypeChanged:
					diffType = "changed"
				case graveler.DiffTypeConflict:
					diffType = "conflicting change"
				case graveler.DiffTypeRemoved:
					diffType = "deleted"
				default:
					diffType = "unknown diff state"
				}
				results = append(results, [2]string{diffType, current.Key.String()})
			}
			w.WriteArray(len(results))
			for _, r := range results {
				w.WriteArray(2)
				w.WriteSimpleString(r[1])
				w.WriteSimpleString(r[0])
			}
		},
	})
	router.Add(&resp.Command{
		Name:     "tag create",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			tag := string(args[0])
			err = store.CreateTag(request.Context(), repo, graveler.TagID(tag), graveler.CommitID(ref))
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			// checkout
			clientName := request.Context().Value(redisClientNameContextKey).(string)
			refStartsAt := strings.LastIndex(clientName, "/") + 1
			clientName = clientName[0:refStartsAt] + tag
			request.SetContext(context.WithValue(request.Context(), redisClientNameContextKey, clientName))
			w.WriteOK()
		},
	})
	router.Add(&resp.Command{
		Name:     "branch list",
		Validate: resp.NoArgs(),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, _, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			branches := make([]string, 0)
			iter, err := store.ListBranches(request.Context(), repo)
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			for iter.Next() {
				branches = append(branches, string(iter.Value().BranchID))
			}
			w.WriteArray(len(branches))
			for _, b := range branches {
				w.WriteSimpleString(b)
			}
		},
	})
	router.Add(&resp.Command{
		Name:     "commit",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			message := string(args[0])
			now := time.Now().Unix()
			_, err = store.Commit(request.Context(), repo, graveler.BranchID(ref), graveler.CommitParams{
				Committer: "redis",
				Message:   message,
				Date:      &now,
			})
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			w.WriteOK()
		},
	})
	router.Add(&resp.Command{
		Name:     "merge",
		Validate: resp.NArgs(1),
		HandlerFn: func(request resp.Request, args [][]byte, w resp.ResponseWriter) {
			repo, ref, err := extractVersion(store, request.Context())
			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			now := time.Now().Unix()
			_, err = store.Merge(request.Context(), repo, graveler.BranchID(ref), graveler.Ref(args[0]), graveler.CommitParams{
				Committer: "redis",
				Message:   fmt.Sprintf("merge %s into %s", args[0], ref),
				Date:      &now,
			}, graveler.MergeStrategySrcWinsStr)

			if err != nil {
				w.WriteSimpleError(resp.ErrorPrefixGeneric, err.Error())
				return
			}
			w.WriteOK()
		},
	})
	return router
}
