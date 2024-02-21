package webdav

import (
	"context"
	"errors"
	"os"
	"path"

	"github.com/go-openapi/swag"

	"github.com/treeverse/lakefs/pkg/uri"

	"github.com/treeverse/lakefs/pkg/api/apigen"
	dav "golang.org/x/net/webdav"
)

type lakeFSROFilesystem struct {
	server apigen.ClientWithResponsesInterface
	root   *uri.URI
}

func (fs *lakeFSROFilesystem) uriFor(name string) *uri.URI {
	p := name
	if fs.root.GetPath() != "" {
		p = path.Join(fs.root.GetPath(), name)
	}
	return &uri.URI{
		Repository: fs.root.Repository,
		Ref:        fs.root.Ref,
		Path:       swag.String(p),
	}
}

func (fs *lakeFSROFilesystem) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return nil
}

func (fs *lakeFSROFilesystem) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (dav.File, error) {
	// try reading file
	fileUri := fs.uriFor(name)
	info, err := getFileInfo(ctx, fs.server, fileUri)
	if errors.Is(err, os.ErrNotExist) {
		// not a file, but perhaps a directory?
		dirInfo, err := getDirInfo(ctx, fs.server, fileUri)
		return &lakeFSROFile{
			server:   fs.server,
			location: fileUri,
			info:     dirInfo,
		}, err
	} else if err != nil {
		// something bad happened!
		return nil, err
	}

	// valid file!
	return &lakeFSROFile{
		server:   fs.server,
		location: fileUri,
		info:     info,
	}, err
}

func (fs *lakeFSROFilesystem) RemoveAll(ctx context.Context, name string) error {
	return nil
}

func (fs *lakeFSROFilesystem) Rename(ctx context.Context, oldName, newName string) error {
	return nil
}

func (fs *lakeFSROFilesystem) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	// try reading file
	fileUri := fs.uriFor(name)
	info, err := getFileInfo(ctx, fs.server, fileUri)
	if errors.Is(err, os.ErrNotExist) {
		// not a file, but perhaps a directory?
		return getDirInfo(ctx, fs.server, fileUri)
	}
	return info, err
}
