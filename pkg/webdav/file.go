package webdav

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/go-openapi/swag"

	"github.com/treeverse/lakefs/pkg/uri"
)

type lakeFSROFile struct {
	server        LakeFSServer
	location      *uri.URI
	fileCache     *fileCache
	metadataCache *metadataCache

	info       *lakeFSFileInfo
	handle     *os.File
	skipHidden bool

	virtualOffset int64
}

func (f *lakeFSROFile) getFile() (*os.File, error) {
	if f.handle != nil {
		return f.handle, nil
	}
	key := f.info.stat.Checksum
	expected := swag.Int64Value(f.info.stat.SizeBytes)
	file, err := f.fileCache.Get(key)
	if err == nil {
		return file, nil
	}
	reader, err := readFile(context.TODO(), f.server, f.location)
	if err != nil {
		return nil, err
	}
	file, err = f.fileCache.Set(f.info.stat.Checksum, reader, expected)
	if err != nil {
		return nil, err
	}
	f.handle = file
	return f.handle, nil
}

func (f *lakeFSROFile) Close() error {
	if f.handle != nil {
		return f.handle.Close()
	}
	return nil
}

func (f *lakeFSROFile) Read(p []byte) (n int, err error) {
	base := path.Base(f.location.GetPath())
	if base != "." && strings.HasPrefix(base, ".") {
		return 0, os.ErrNotExist
	}
	file, err := f.getFile()
	if err != nil {
		return 0, err
	}
	_, err = file.Seek(f.virtualOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}
	n, err = file.Read(p)
	f.virtualOffset += int64(n)
	return n, err
}

func (f *lakeFSROFile) Seek(offset int64, whence int) (int64, error) {
	base := path.Base(f.location.GetPath())
	if base != "." && strings.HasPrefix(base, ".") {
		return 0, os.ErrNotExist
	}
	file, err := f.getFile()
	if err != nil {
		return 0, err
	}
	offset, err = file.Seek(offset, whence)
	if err != nil {
		return 0, err
	}
	f.virtualOffset = offset
	return offset, err
}

func (f *lakeFSROFile) Readdir(count int) ([]fs.FileInfo, error) {
	base := path.Base(f.location.GetPath())
	if base != "." && strings.HasPrefix(base, ".") {
		return nil, os.ErrNotExist
	}
	var listing []*lakeFSFileInfo
	var err error
	if f.skipHidden {
		listing, err = listDirectory(context.TODO(), f.metadataCache, f.server, f.location, count, skipEmpty, skipHidden)
	} else {
		listing, err = listDirectory(context.TODO(), f.metadataCache, f.server, f.location, count, skipEmpty)
	}
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
	base := path.Base(f.location.GetPath())
	if base != "." && strings.HasPrefix(base, ".") {
		return nil, os.ErrNotExist
	}
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
	return 0, nil
}
