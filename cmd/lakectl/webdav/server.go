package webdav

import (
	"context"
	"net/http"
	"os"

	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/uri"

	dav "golang.org/x/net/webdav"
)

func NewServer(location *uri.URI, lakefs apigen.ClientWithResponsesInterface) (http.Handler, error) {
	fs := dav.NewMemFS()
	f, err := fs.OpenFile(context.TODO(), "hello.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	_, _ = f.Write([]byte("hello world!!!\n"))
	_ = f.Close()

	handler := &dav.Handler{
		Prefix:     "/",
		FileSystem: fs,
		LockSystem: dav.NewMemLS(),
		Logger:     nil,
	}

	return handler, nil
}
