package diff

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/treeverse/lakefs/pkg/plugins"
)

type Differ interface {
	Diff(context.Context, TablePaths, S3Creds) ([]*Diff, error)
}

var (
	ErrUninitializedDiffService = errors.New("uninitialized diff service")
)

type TablePaths struct {
	LeftTablePath  string
	RightTablePath string
	BaseTablePath  string
}

type S3Creds struct {
	Key      string
	Secret   string
	Endpoint string
}

type Params struct {
	TablePaths TablePaths
	S3Creds    S3Creds
}

type Entry struct {
	Version          string
	Timestamp        time.Time
	Operation        string
	OperationContent map[string]string
}

type Service struct {
	pluginManager *plugins.Manager[Differ]
}

func NewService(pm *plugins.Manager[Differ]) *Service {
	return &Service{
		pluginManager: pm,
	}
}

func (s *Service) RunDiff(ctx context.Context, diffName string, diffParams Params) ([]Entry, error) {
	if s == nil {
		return nil, ErrUninitializedDiffService
	}
	d, closeClient, err := s.pluginManager.LoadPluginClient(diffName)
	if err != nil {
		return nil, err
	}
	defer closeClient()
	diffs, err := (*d).Diff(ctx, diffParams.TablePaths, diffParams.S3Creds)
	if err != nil {
		return nil, err
	}
	return buildEntries(diffs), nil
}

func buildEntries(diffs []*Diff) []Entry {
	result := make([]Entry, len(diffs))
	for _, diff := range diffs {
		result = append(result, Entry{
			Version:          strconv.Itoa(int(diff.Version)),
			Timestamp:        diff.Timestamp.AsTime(),
			Operation:        diff.Description,
			OperationContent: diff.Content,
		})
	}
	return result
}