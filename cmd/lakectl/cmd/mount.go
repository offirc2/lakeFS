package cmd

import (
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/treeverse/lakefs/cmd/lakectl/webdav"
)

// logCmd represents the log command
var mountCmd = &cobra.Command{
	Use:   "mount <path URI> <location>",
	Short: "mount a path to a local directory",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		loc := args[1]
		server, err := webdav.NewServer(MustParsePathURI("path URI", args[0]), getClient())
		if err != nil {
			DieErr(err)
		}

		listener, err := net.Listen("tcp4", "127.0.0.1:0")
		if err != nil {
			DieErr(err)
		}
		addr, ok := listener.Addr().(*net.TCPAddr)
		if !ok {
			Die("could not read listener port", 1)
		}
		fmt.Printf("listening on 127.0.0.1:%d\n", addr.Port)
		done := make(chan struct{})

		go func() {
			if err := http.Serve(listener, server); err != nil {
				DieErr(err)
			}
			done <- struct{}{}
		}()

		// mount
		if err := webdav.Mount(addr, loc); err != nil {
			DieErr(err)
		}
		fmt.Printf("mounted!")
		<-done
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(mountCmd)

}
