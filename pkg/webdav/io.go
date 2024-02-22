package webdav

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/uri"
)

const (
	PathDelimiter = "/"
	MaxPageSize   = 1000
)

var (
	ErrLakeFSError = errors.New("lakeFS error")
)

type filterFn func(parent, path string) bool

func skipHidden(parent, p string) bool {
	base := path.Base(p)
	isHidden := base != "." && strings.HasPrefix(base, ".")
	return !isHidden
}

func skipEmpty(parent, p string) bool {
	if len(p) < len(parent) {
		return false
	}
	return true
}

func listDirectory(ctx context.Context, cache *metadataCache, server apigen.ClientWithResponsesInterface, location *uri.URI, amount int, filters ...filterFn) ([]*lakeFSFileInfo, error) {
	dirPath := location.GetPath()

	cacheKey := metadataCacheKey{
		repo: location.Repository,
		ref:  location.Ref,
		path: dirPath,
	}

	if dirPath != "" && !strings.HasSuffix(dirPath, PathDelimiter) {
		dirPath += PathDelimiter
	}

	listingAmount := amount
	if amount <= 0 || amount >= MaxPageSize || len(filters) > 0 {
		listingAmount = MaxPageSize
	}

	if amount <= 0 {
		// we currently only cache full listings
		if cachedListing, hit := cache.getListing(cacheKey); hit {
			return cachedListing, nil
		}
	}

	prefix := apigen.PaginationPrefix(dirPath)
	listingAmountParam := apigen.PaginationAmount(listingAmount)
	delimiter := apigen.PaginationDelimiter(PathDelimiter)

	hasMore := true
	nextOffset := ""
	results := make([]*lakeFSFileInfo, 0)
	for hasMore && (amount <= 0 || len(results) < amount) {
		after := apigen.PaginationAfter(nextOffset)
		response, err := server.ListObjectsWithResponse(ctx, location.Repository, location.Ref, &apigen.ListObjectsParams{
			Amount:    &listingAmountParam,
			Delimiter: &delimiter,
			Prefix:    &prefix,
			After:     &after,
		})
		if err != nil {
			return nil, err
		}
		if response.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("%w: HTTP %d", ErrLakeFSError, response.StatusCode())
		}

		for _, result := range response.JSON200.Results {
			p := result.Path
			if strings.HasSuffix(p, PathDelimiter) {
				p = p[0 : len(p)-1]
			}
			passed := true
			for _, filter := range filters {
				if !filter(dirPath, p) {
					passed = false
				}
			}
			if !passed {
				continue
			}
			results = append(results, &lakeFSFileInfo{
				location: &uri.URI{
					Repository: location.Repository,
					Ref:        location.Ref,
					Path:       swag.String(p),
				},
				dir:  result.PathType == "common_prefix",
				stat: &result,
			})
			if amount > 0 && len(results) >= amount {
				break
			}
		}
		hasMore = response.JSON200.Pagination.HasMore
		nextOffset = response.JSON200.Pagination.NextOffset
	}

	if amount <= 0 {
		// only cache full results
		cache.setListing(cacheKey, results)
	}
	return results, nil
}

func getDirInfo(ctx context.Context, cache *metadataCache, server apigen.ClientWithResponsesInterface, location *uri.URI) (*lakeFSFileInfo, error) {
	cacheKey := metadataCacheKey{
		repo: location.Repository,
		ref:  location.Ref,
		path: location.GetPath(),
	}
	_, hit := cache.getListing(cacheKey)
	if hit {
		return &lakeFSFileInfo{
			location: location,
			dir:      true,
		}, nil // cache hit!
	}
	listing, err := listDirectory(ctx, cache, server, location, 1)
	if err != nil {
		return nil, err
	}
	if len(listing) == 0 {
		return nil, os.ErrNotExist
	}
	return &lakeFSFileInfo{
		location: location,
		dir:      true,
	}, nil
}

func getFileInfo(ctx context.Context, cache *metadataCache, server apigen.ClientWithResponsesInterface, location *uri.URI) (*lakeFSFileInfo, error) {
	if location.GetPath() == "" {
		// lakeFS can't stat an empty path
		return nil, os.ErrNotExist
	}
	cacheKey := metadataCacheKey{
		repo: location.Repository,
		ref:  location.Ref,
		path: location.GetPath(),
	}
	cached, hit := cache.getObject(cacheKey)
	if hit {
		if cached == nil {
			return nil, os.ErrNotExist
		} else {
			return cached, nil
		}
	}
	response, err := server.StatObjectWithResponse(
		ctx, location.Repository, location.Ref, &apigen.StatObjectParams{Path: location.GetPath()})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() == http.StatusOK {
		info := &lakeFSFileInfo{
			location: location,
			dir:      false,
			stat:     response.JSON200,
		}
		cache.setObject(cacheKey, info)
		return info, nil
	} else if response.StatusCode() == http.StatusNotFound {
		cache.setObject(cacheKey, nil)
		return nil, os.ErrNotExist
	}
	return nil, fmt.Errorf("%w: HTTP %d", ErrLakeFSError, response.StatusCode())
}

func readFile(ctx context.Context, server apigen.ClientInterface, location *uri.URI) (io.ReadCloser, error) {
	response, err := server.GetObject(ctx, location.Repository, location.Ref, &apigen.GetObjectParams{
		Path:    location.GetPath(),
		Presign: swag.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 299 {
		return nil, fmt.Errorf("%w: HTTP %d", ErrLakeFSError, response.StatusCode)
	}
	return response.Body, nil
}
