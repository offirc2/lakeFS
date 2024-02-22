package webdav

import (
	"github.com/treeverse/lakefs/pkg/uri"
	"gopkg.in/yaml.v3"
)

type MountFileEntry struct {
	LocalPath  string   `yaml:"local_path"`
	RemotePath *uri.URI `yaml:"remote_path"`
	Head       string   `yaml:"head,omitempty"`
	Mode       string   `yaml:"mode"`
}

type MountFile struct {
	Mounts []MountFileEntry `yaml:"mounts"`
}

func ReadMountFile(data []byte) (*MountFile, error) {
	mounts := &MountFile{}
	err := yaml.Unmarshal(data, mounts)
	return mounts, err
}

func (m *MountFile) Write() ([]byte, error) {
	return yaml.Marshal(m)
}
