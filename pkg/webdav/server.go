package webdav

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"

	dav "golang.org/x/net/webdav"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"

	"github.com/treeverse/lakefs/pkg/uri"
)

type MountHTTPError struct {
	Error string `json:"error"`
}

type MountInfo struct {
	Remote    *uri.URI `json:"remote"`
	LocalPath string   `json:"local_path"`
	Mode      string   `json:"mode"`
}

type MountRegistry struct {
	mounts map[string]*MountInfo
	lock   sync.RWMutex
}

func (r *MountRegistry) ListMounts() ([]*MountInfo, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	mounts := make([]*MountInfo, len(r.mounts))
	i := 0
	for _, v := range r.mounts {
		mounts[i] = v
		i++
	}
	sort.Slice(mounts, func(i, j int) bool {
		return mounts[i].LocalPath < mounts[j].LocalPath
	})
	return mounts, nil
}

func (r *MountRegistry) Unmount(localPath string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.mounts, localPath)
	return nil
}

func (r *MountRegistry) Mount(mountInfo *MountInfo) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.mounts[mountInfo.LocalPath] = mountInfo
	return nil
}

func writeError(writer http.ResponseWriter, statusCode int, err error) {
	_ = writeJSON(writer, statusCode, MountHTTPError{Error: err.Error()})
}
func writeJSON(writer http.ResponseWriter, statusCode int, data any) error {
	serialized, _ := json.Marshal(data)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_, err := writer.Write(serialized)
	return err
}

func readJSON(request *http.Request, target any) error {
	return json.NewDecoder(request.Body).Decode(target)
}

func NewServer(addr, cacheDirectory string, server LakeFSServer) (*http.Server, error) {
	registry := &MountRegistry{
		mounts: make(map[string]*MountInfo),
		lock:   sync.RWMutex{},
	}
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	webServer := &http.Server{
		Addr: addr,
	}

	// setup WebDAV Handler (read-only)
	readOnlyFS := &lakeFSROFilesystem{
		server:        server,
		fileCache:     newFileCache(cacheDirectory),
		metadataCache: newMetadataCache(),
		skipHidden:    true,
	}

	router.Delete("/mounts", func(writer http.ResponseWriter, request *http.Request) {
		// delete mount
		localPath := request.URL.Query().Get("path")
		_ = registry.Unmount(localPath)
		writer.WriteHeader(http.StatusOK)
	})

	router.Post("/mounts", func(writer http.ResponseWriter, request *http.Request) {
		// create new mount
		mount := &MountInfo{}
		err := readJSON(request, mount)
		if err != nil {
			writeError(writer, http.StatusBadRequest, err)
			return
		}

		err = registry.Mount(mount)
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}

		_ = writeJSON(writer, http.StatusCreated, mount)
	})

	router.Get("/mounts", func(writer http.ResponseWriter, request *http.Request) {
		// list mounts
		mounts, err := registry.ListMounts()
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err)
			return
		}
		_ = writeJSON(writer, http.StatusOK, mounts)
	})

	router.Post("/terminate", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		go func() {
			_ = webServer.Shutdown(context.TODO())
		}()
	})

	// webdav goes here
	const readOnlyMount = "/wd/read-only"
	webDavHandler := &dav.Handler{
		Prefix:     readOnlyMount,
		FileSystem: readOnlyFS,
		LockSystem: dav.NewMemLS(),
		Logger: func(request *http.Request, err error) {
			//fmt.Printf("WebDav: got request for %s %s\n", request.Method, request.URL.Path)
		},
	}

	webServer.Handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.HasPrefix(request.URL.Path, readOnlyMount) {
			webDavHandler.ServeHTTP(writer, request)
			return
		}

		// otherwise
		router.ServeHTTP(writer, request)
	})

	return webServer, nil
}

func IsServerRunning(addr string) (bool, error) {
	ln, err := net.Listen("tcp4", addr)
	defer func() {
		if ln != nil {
			_ = ln.Close()
		}
	}()
	if err != nil && isErrorAddressAlreadyInUse(err) {
		return true, nil
	}
	return false, err
}

// taken from: https://stackoverflow.com/a/65865898
func isErrorAddressAlreadyInUse(err error) bool {
	var eOsSyscall *os.SyscallError
	if !errors.As(err, &eOsSyscall) {
		return false
	}
	var errErrno syscall.Errno // doesn't need a "*" (ptr) because it's already a ptr (uintptr)
	if !errors.As(eOsSyscall, &errErrno) {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}
	return false
}
