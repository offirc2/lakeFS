package webdav

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type metadataCacheKey struct {
	repo, ref, path string
}

type metadataCache struct {
	objects  map[metadataCacheKey]*lakeFSFileInfo
	listings map[metadataCacheKey][]*lakeFSFileInfo
	lock     *sync.RWMutex
}

func newMetadataCache() *metadataCache {
	return &metadataCache{
		objects:  make(map[metadataCacheKey]*lakeFSFileInfo),
		listings: make(map[metadataCacheKey][]*lakeFSFileInfo),
		lock:     &sync.RWMutex{},
	}
}

func (c *metadataCache) getObject(key metadataCacheKey) (*lakeFSFileInfo, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	info, ok := c.objects[key]
	return info, ok
}

func (c *metadataCache) setObject(key metadataCacheKey, info *lakeFSFileInfo) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.objects[key] = info
}

func (c *metadataCache) getListing(key metadataCacheKey) ([]*lakeFSFileInfo, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	infos, ok := c.listings[key]
	return infos, ok
}

func (c *metadataCache) setListing(key metadataCacheKey, infos []*lakeFSFileInfo) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.listings[key] = infos
}

type fileCache struct {
	dir string
}

func newFileCache(dir string) *fileCache {
	return &fileCache{dir: dir}
}

func (c *fileCache) Get(key string) (*os.File, error) {
	path := filepath.Join(c.dir, key)
	return os.Open(path)
}

func (c *fileCache) Set(key string, content io.ReadCloser, expected int64) (*os.File, error) {
	path := filepath.Join(c.dir, fmt.Sprintf("%s-w", key))
	out, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	n, err := io.Copy(out, content)
	if err != nil {
		// we now have a bad file on our hands
		_ = out.Close()
		_ = os.Remove(path)
		return nil, err
	}
	if expected > 0 && n != expected {
		return nil, os.ErrInvalid
	}

	err = out.Close()
	if err != nil {
		return nil, err
	}
	// make available
	err = os.Rename(path, filepath.Join(c.dir, key))
	if err != nil {
		return nil, err
	}
	f, err := c.Get(key)
	return f, err
}
