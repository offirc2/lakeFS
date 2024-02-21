package webdav

import (
	"bytes"
	"context"
	"io/fs"
	"os"

	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/uri"
)

type lakeFSROFile struct {
	server   apigen.ClientWithResponsesInterface
	location *uri.URI

	info *lakeFSFileInfo
}

func (f *lakeFSROFile) Close() error {
	//TODO implement me
	panic("implement me")
}

func (f *lakeFSROFile) Read(p []byte) (n int, err error) {
	data, err := readFile(context.TODO(), f.server, f.location)
	if err != nil {
		return len(data), err
	}
	buf := bytes.NewBuffer(nil)
	buf.Write(data)
	return buf.Read(p) // LOL this is such a stupid impl
}

func (f *lakeFSROFile) Seek(offset int64, whence int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (f *lakeFSROFile) Readdir(count int) ([]fs.FileInfo, error) {
	listing, err := listDirectory(context.TODO(), f.server, f.location, count)
	if err != nil {
		return nil, err
	}
	results := make([]fs.FileInfo, len(listing))
	for i, entry := range listing {
		results[i] = entry
	}
	return results, nil
}

func (f *lakeFSROFile) Stat() (fs.FileInfo, error) {
	if f.info != nil {
		return f.info, nil
	}

	if f.info.IsDir() {
		return &lakeFSFileInfo{
			location: f.location,
			dir:      true,
		}, nil
	}

	return nil, os.ErrNotExist
}

func (f *lakeFSROFile) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}
