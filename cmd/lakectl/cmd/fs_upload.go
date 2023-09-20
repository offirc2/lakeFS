package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/api/apiutil"
	"github.com/treeverse/lakefs/pkg/api/helpers"
	"github.com/treeverse/lakefs/pkg/uri"
)

const fsRecursiveTemplate = `Files: {{.Count}}
Total Size: {{.Bytes}} bytes
Human Total Size: {{.Bytes|human_bytes}}
`

var fsUploadCmd = &cobra.Command{
	Use:               "upload <path uri>",
	Short:             "Upload a local file to the specified URI",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: ValidArgsRepository,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		pathURI := MustParsePathURI("path", args[0])
		flagSet := cmd.Flags()
		source := Must(flagSet.GetString("source"))
		recursive := Must(flagSet.GetBool("recursive"))
		preSignMode := Must(flagSet.GetBool("pre-sign"))
		contentType := Must(flagSet.GetString("content-type"))

		ctx := cmd.Context()
		if !recursive {
			if pathURI.GetPath() == "" {
				Die("target path is not a valid URI", 1)
			}
			stat, err := upload(ctx, client, source, pathURI, contentType, preSignMode)
			if err != nil {
				DieErr(err)
			}
			Write(fsStatTemplate, stat)
			return
		}

		// copy recursively
		var totals struct {
			Bytes int64
			Count int64
		}
		err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("traverse %s: %w", path, err)
			}
			if info.IsDir() {
				return nil
			}
			relPath := strings.TrimPrefix(path, source)
			uri := *pathURI
			p := filepath.ToSlash(filepath.Join(*uri.Path, relPath))
			uri.Path = &p
			stat, err := upload(ctx, client, path, &uri, contentType, preSignMode)
			if err != nil {
				return fmt.Errorf("upload %s: %w", path, err)
			}
			if stat.SizeBytes != nil {
				totals.Bytes += *stat.SizeBytes
			}
			totals.Count++
			return nil
		})
		if err != nil {
			DieErr(err)
		}
		Write(fsRecursiveTemplate, totals)
	},
}

func upload(ctx context.Context, client apigen.ClientWithResponsesInterface, sourcePathname string, destURI *uri.URI, contentType string, preSign bool) (*apigen.ObjectStats, error) {
	fp := Must(OpenByPath(sourcePathname))
	defer func() {
		_ = fp.Close()
	}()
	objectPath := apiutil.Value(destURI.Path)
	if preSign {
		return helpers.ClientUploadPreSign(ctx, client, destURI.Repository, destURI.Ref, objectPath, nil, contentType, fp)
	}
	return helpers.ClientUpload(ctx, client, destURI.Repository, destURI.Ref, objectPath, nil, contentType, fp)
}

//nolint:gochecknoinits
func init() {
	fsUploadCmd.Flags().StringP("source", "s", "", "local file to upload, or \"-\" for stdin")
	fsUploadCmd.Flags().BoolP("recursive", "r", false, "recursively copy all files under local source")
	_ = fsUploadCmd.MarkFlagRequired("source")
	fsUploadCmd.Flags().StringP("content-type", "", "", "MIME type of contents")
	fsUploadCmd.Flags().Bool("pre-sign", false, "Use pre-sign link to access the data")

	fsCmd.AddCommand(fsUploadCmd)
}
