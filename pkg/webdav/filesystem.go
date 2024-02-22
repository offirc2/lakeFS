package webdav

import (
	"context"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/uri"

	dav "golang.org/x/net/webdav"
)

type LakeFSServer interface {
	apigen.ClientWithResponsesInterface
	apigen.ClientInterface
}

type lakeFSROFilesystem struct {
	server        LakeFSServer
	fileCache     *fileCache
	metadataCache *metadataCache
	skipHidden    bool
}

func UriFor(name string) (*uri.URI, error) {
	if strings.HasPrefix(name, "/") {
		name = name[1:] // remove leading slash
	}
	parts := strings.SplitN(name, "/", 4)
	if len(parts) == 4 {
		// nonce, repository, reference, path
		return &uri.URI{Repository: parts[1], Ref: parts[2], Path: swag.String(parts[3])}, nil
	} else if len(parts) == 3 && parts[2] != "" {
		// nonce, repository, reference
		return &uri.URI{Repository: parts[1], Ref: parts[2], Path: swag.String("")}, nil
	} else {
		// invalid?
		return nil, os.ErrNotExist
	}
}

func (fs *lakeFSROFilesystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return nil
}

func (fs *lakeFSROFilesystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (dav.File, error) {
	// try reading file
	fileUri, err := UriFor(name)
	if err != nil {
		return nil, err
	}
	info, err := getFileInfo(ctx, fs.metadataCache, fs.server, fileUri)
	if errors.Is(err, os.ErrNotExist) {
		// not a file, but perhaps a directory?
		dirInfo, err := getDirInfo(ctx, fs.metadataCache, fs.server, fileUri)
		return &lakeFSROFile{
			server:        fs.server,
			location:      fileUri,
			info:          dirInfo,
			fileCache:     fs.fileCache,
			metadataCache: fs.metadataCache,
			skipHidden:    fs.skipHidden,
		}, err
	} else if err != nil {
		// something bad happened!
		return nil, err
	}

	// valid file!
	return &lakeFSROFile{
		server:        fs.server,
		location:      fileUri,
		info:          info,
		fileCache:     fs.fileCache,
		metadataCache: fs.metadataCache,
		skipHidden:    fs.skipHidden,
	}, err
}

func (fs *lakeFSROFilesystem) RemoveAll(ctx context.Context, name string) error {
	return nil
}

func (fs *lakeFSROFilesystem) Rename(ctx context.Context, oldName, newName string) error {
	return nil
}

func (fs *lakeFSROFilesystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	base := path.Base(name)
	if fs.skipHidden && base != "." && strings.HasPrefix(base, ".") {
		return nil, os.ErrNotExist
	}
	// try reading file
	fileUri, err := UriFor(name)
	if err != nil {
		return nil, err
	}
	info, err := getFileInfo(ctx, fs.metadataCache, fs.server, fileUri)
	if errors.Is(err, os.ErrNotExist) {
		// not a file, but perhaps a directory?
		return getDirInfo(ctx, fs.metadataCache, fs.server, fileUri)
	}
	return info, err
}
