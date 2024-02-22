package webdav

import (
	"context"
	"io/fs"
	"os"
	"path"
	"time"

	"golang.org/x/net/webdav"

	"github.com/go-openapi/swag"
	"github.com/treeverse/lakefs/pkg/uri"

	"github.com/treeverse/lakefs/pkg/api/apigen"
)

type lakeFSFileInfo struct {
	location *uri.URI
	dir      bool
	stat     *apigen.ObjectStats
}

func (fi *lakeFSFileInfo) Name() string {
	return path.Base(fi.location.GetPath())
}

func (fi *lakeFSFileInfo) Size() int64 {
	if fi.dir || fi.stat == nil {
		return 0
	}
	return swag.Int64Value(fi.stat.SizeBytes)
}

func (fi *lakeFSFileInfo) Mode() fs.FileMode {
	if fi.dir {
		return os.ModeDir
	}
	return os.ModePerm
}

func (fi *lakeFSFileInfo) ModTime() time.Time {
	if fi.dir || fi.stat == nil {
		return time.Now()
	}
	return time.Unix(fi.stat.Mtime, 0).UTC().Local()
}

func (fi *lakeFSFileInfo) IsDir() bool {
	return fi.dir
}

func (fi *lakeFSFileInfo) ETag(ctx context.Context) (string, error) {
	if fi.stat != nil {
		return fi.stat.Checksum, nil
	}
	return "", webdav.ErrNotImplemented
}

func (fi *lakeFSFileInfo) Sys() any {
	webdav.NewMemFS()
	return fi.stat
}
