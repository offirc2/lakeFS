package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/treeverse/lakefs/pkg/fileutil"
	"github.com/treeverse/lakefs/pkg/uri"
	"github.com/treeverse/lakefs/pkg/webdav"
)

const (
	MountServerBindAddress = "127.0.0.1:6363"
)

func runMountFreeze() {
	if !isMountServerAlive() {
		return
	}
	restClient := webdav.NewMountServerRestClient(MountServerBindAddress)
	mounts := Must(restClient.GetMounts())
	fstab := &webdav.MountFile{Mounts: make([]webdav.MountFileEntry, 0)}
	for _, mount := range mounts {
		localPath := mount.LocalPath // expected to be full path
		cwd := Must(os.Getwd())
		rel, err := filepath.Rel(cwd, localPath)
		rel = filepath.ToSlash(rel)
		if err != nil || strings.HasPrefix(rel, "..") {
			continue // not a sub path!
		}

		fstab.Mounts = append(fstab.Mounts,
			webdav.MountFileEntry{
				LocalPath:  rel,
				RemotePath: mount.Remote,
			},
		)
	}
	data := Must(fstab.Write())
	fmt.Printf("%s\n", data)
}

func runMountFrom(from string) {
	inputFile := Must(os.Open(from))
	data := Must(io.ReadAll(inputFile))
	mounts := Must(webdav.ReadMountFile(data))
	for _, mnt := range mounts.Mounts {
		runMount(mnt.RemotePath, mnt.LocalPath)
	}
}

func runMountConfig(key, value string) {
	fmt.Printf("config pair: %s = %s!\n", key, value)
	return
}

func isMountServerAlive() bool {
	return Must(webdav.IsServerRunning(MountServerBindAddress))
}

func ensureMountServerRunning() {
	if !isMountServerAlive() {
		// no server, let's spawn one!
		pid := Must(webdav.Daemonize("mount", "--server"))
		fmt.Printf("started mount server with pid %d\n", pid)
		// wait for it to be up
		attempts := 3
		up := false
		for i := 0; i < attempts; i++ {
			if isMountServerAlive() {
				up = true
				break
			}
			time.Sleep(time.Second)
		}
		if !up {
			Die("could not spin up local mount server", 1)
		}
	}
}

func runUmount(local string) {
	absolutePath := Must(filepath.Abs(local))
	if !isMountServerAlive() {
		return
	}
	if err := webdav.Umount(absolutePath); err != nil {
		DieErr(err)
	}

	restClient := webdav.NewMountServerRestClient(MountServerBindAddress)
	if err := restClient.Unmount(absolutePath); err != nil {
		DieErr(err)
	}

	if len(Must(restClient.GetMounts())) == 0 {
		// no mounts left, stop server
		if err := restClient.TerminateServer(); err != nil {
			DieErr(err)
		}
	}

}

func getGitRoot(path string) (string, error) {
	gitRoot, err := fileutil.FindInParents(path, ".git")
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return gitRoot, nil
}

func runMount(remote *uri.URI, local string) {
	// let's see if there's a running server:
	ensureMountServerRunning()

	// make an absolute path
	absolutePath := Must(filepath.Abs(local))
	gitRoot := Must(getGitRoot(absolutePath))
	fmt.Printf("git root: %s\n", gitRoot)

	// now that we have a mount server, let's mount it!
	restClient := webdav.NewMountServerRestClient(MountServerBindAddress)
	// ensure dir
	if !Must(fileutil.IsDir(absolutePath)) {
		err := os.MkdirAll(absolutePath, 0755)
		if err != nil {
			DieErr(err)
		}
	}
	err := webdav.Mount(restClient.GetWebdavURL("read-only", remote), absolutePath)
	if err != nil {
		DieErr(err)
	}

	err = restClient.RegisterMount(remote, absolutePath, "read-only")
	if err != nil {
		DieErr(err)
	}
}

func runMountServer() {
	client := getClient()
	cacheDir := os.Getenv("LAKECTL_MOUNT_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), "lakefs-mount-cache")
	}
	if !Must(fileutil.IsDir(cacheDir)) {
		err := os.MkdirAll(cacheDir, 0755)
		if err != nil {
			DieErr(err)
		}
	}
	if err := webdav.RunServer(MountServerBindAddress, cacheDir, client); err != nil {
		DieErr(err)
	}
}

var mountCmd = &cobra.Command{
	Use:   "mount <path URI> <location>",
	Short: "mount a path to a local directory",
	Args: func(cmd *cobra.Command, args []string) error {
		freeze := Must(cmd.Flags().GetBool("freeze"))
		if freeze {
			if len(args) > 0 {
				return fmt.Errorf("freeze takes no positional arguments")
			}
			return nil
		}

		server := Must(cmd.Flags().GetBool("server"))
		if server {
			if len(args) > 0 {
				return fmt.Errorf("freeze takes no positional arguments")
			}
			return nil
		}

		from := Must(cmd.Flags().GetString("from"))
		if from != "" {
			// validate it's a file
			if !Must(fileutil.FileExists(from)) {
				return fmt.Errorf("file not found: %s\n", from)
			}
			return nil
		}

		configPair := Must(cmd.Flags().GetString("config"))
		if configPair != "" {
			_, _, isKV := strings.Cut(configPair, "=")
			if !isKV {
				return fmt.Errorf("usage: --config <key> <value>\n")
			}
			return nil
		}

		// regular mount!
		return cobra.ExactArgs(2)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		freeze := Must(cmd.Flags().GetBool("freeze"))
		server := Must(cmd.Flags().GetBool("server"))
		from := Must(cmd.Flags().GetString("from"))
		conf := Must(cmd.Flags().GetString("config"))
		if server {
			runMountServer()
		} else if freeze {
			runMountFreeze()
		} else if from != "" {
			runMountFrom(from)
		} else if conf != "" {
			key, value, _ := strings.Cut(conf, "=")
			runMountConfig(key, value)
		} else {
			runMount(MustParsePathURI("remote", args[0]), args[1])
		}
	},
}

var umountCmd = &cobra.Command{
	Use:   "umount <location>",
	Short: "unmount a path",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runUmount(args[0])
	},
}

//nolint:gochecknoinits
func init() {
	mountCmd.Flags().Bool("freeze", false, "print currently mounted paths")
	mountCmd.Flags().String("from", "", "mount paths as listed in file")
	_ = mountCmd.MarkFlagFilename("from")
	mountCmd.Flags().String("config", "", "set configuration parameters for mounts \"key=value\"")
	mountCmd.Flags().Bool("server", false, "")
	_ = mountCmd.Flags().MarkHidden("server") // only used internally for daemonization

	mountCmd.MarkFlagsMutuallyExclusive("freeze", "from", "config", "server")
	rootCmd.AddCommand(mountCmd)

	rootCmd.AddCommand(umountCmd)
}
