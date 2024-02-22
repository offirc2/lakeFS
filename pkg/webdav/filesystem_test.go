package webdav_test

import (
	"os"
	"testing"

	"github.com/treeverse/lakefs/pkg/webdav"

	"github.com/treeverse/lakefs/pkg/uri"
)

func TestURIFor(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    *uri.URI
		ExpectedErr error
	}{
		{Input: "nonce/a/b/c/d/e", Expected: uri.Must(uri.Parse("lakefs://a/b/c/d/e"))},
		{Input: "/nonce/a/b/c/d/e", Expected: uri.Must(uri.Parse("lakefs://a/b/c/d/e"))},
		{Input: "nonce/a/b", Expected: uri.Must(uri.Parse("lakefs://a/b"))},
		{Input: "nonce/a/b/", Expected: uri.Must(uri.Parse("lakefs://a/b"))},
		{Input: "nonce/a/", ExpectedErr: os.ErrNotExist},
		{Input: "nonce/a", ExpectedErr: os.ErrNotExist},
		{Input: "nonce", ExpectedErr: os.ErrNotExist},
		{Input: "/nonce/", ExpectedErr: os.ErrNotExist},
	}
	for _, cas := range cases {
		t.Run(cas.Input, func(t *testing.T) {
			got, err := webdav.UriFor(cas.Input)
			if err != cas.ExpectedErr {
				t.Errorf("unxpected error for case \"%s\": got \"%s\", expected: \"%s\"",
					cas.Input, err, cas.ExpectedErr)
			}
			if err != nil {
				return // no need to check struct if got an error
			}
			if got.Repository != cas.Expected.Repository {
				t.Errorf("unxpected repository for case \"%s\": got \"%s\", expected: \"%s\"",
					cas.Input, got.Repository, cas.Expected.Repository)
			}
			if got.Ref != cas.Expected.Ref {
				t.Errorf("unxpected ref for case %s: got %s, expected: %s",
					cas.Input, got.Ref, cas.Expected.Ref)
			}
			if got.GetPath() != cas.Expected.GetPath() {
				t.Errorf("unxpected path for case %s: got %s, expected: %s",
					cas.Input, got.GetPath(), cas.Expected.GetPath())
			}
		})
	}
}
