package webdav

import (
	"github.com/treeverse/lakefs/pkg/api/apigen"
	"github.com/treeverse/lakefs/pkg/uri"
)

type lakeFSROPrefix struct {
	server   apigen.ClientWithResponsesInterface
	location *uri.URI
}
