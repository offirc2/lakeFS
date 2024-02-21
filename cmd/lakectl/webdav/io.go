package webdav

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
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

func listDirectory(ctx context.Context, server apigen.ClientWithResponsesInterface, location *uri.URI, amount int) ([]*lakeFSFileInfo, error) {
	// how can we tell if there's a dir? only by listing "<name>/"
	path := location.GetPath()
	if strings.HasSuffix(path, PathDelimiter) {
		path += PathDelimiter
	}

	listingAmount := amount
	if amount == -1 || amount >= MaxPageSize {
		listingAmount = MaxPageSize
	}

	prefix := apigen.PaginationPrefix(path)
	listingAmountParam := apigen.PaginationAmount(listingAmount)
	delimiter := apigen.PaginationDelimiter(PathDelimiter)

	hasMore := true
	nextOffset := ""
	results := make([]*lakeFSFileInfo, 0)
	for hasMore && (amount == -1 || len(results) < amount) {
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
			results = append(results, &lakeFSFileInfo{
				location: &uri.URI{
					Repository: location.Repository,
					Ref:        location.Ref,
					Path:       swag.String(p),
				},
				dir:  result.PathType == "common_prefix",
				stat: &result,
			})
			if len(results) >= amount {
				break
			}
		}
		hasMore = response.JSON200.Pagination.HasMore
		nextOffset = response.JSON200.Pagination.NextOffset
	}

	return results, nil
}

func getDirInfo(ctx context.Context, server apigen.ClientWithResponsesInterface, location *uri.URI) (*lakeFSFileInfo, error) {
	listing, err := listDirectory(ctx, server, location, 1)
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

func getFileInfo(ctx context.Context, server apigen.ClientWithResponsesInterface, location *uri.URI) (*lakeFSFileInfo, error) {
	response, err := server.StatObjectWithResponse(
		ctx, location.Repository, location.Ref, &apigen.StatObjectParams{Path: location.GetPath()})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() == http.StatusOK {
		return &lakeFSFileInfo{
			location: location,
			dir:      false,
			stat:     response.JSON200,
		}, nil
	} else if response.StatusCode() == http.StatusNotFound {
		return nil, os.ErrNotExist
	}
	return nil, fmt.Errorf("%w: HTTP %d", ErrLakeFSError, response.StatusCode())
}

func readFile(ctx context.Context, server apigen.ClientWithResponsesInterface, location *uri.URI) ([]byte, error) {
	response, err := server.GetObjectWithResponse(ctx, location.Repository, location.Ref, &apigen.GetObjectParams{
		Path:    location.GetPath(),
		Presign: swag.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	if response.StatusCode() > 299 {
		return nil, fmt.Errorf("%w: HTTP %d", ErrLakeFSError, response.StatusCode())
	}
	return response.Body, nil
}
