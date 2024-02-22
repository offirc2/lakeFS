package webdav

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// fork crete a new process
func fork(args []string) (int, error) {
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = os.Environ()
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.ExtraFiles = nil
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	pid := cmd.Process.Pid
	// release
	if err := cmd.Process.Release(); err != nil {
		return pid, err
	}
	return pid, nil
}

func Daemonize(cmd ...string) (int, error) {
	return fork(cmd)
}

func RunServer(addr, cacheDir string, server LakeFSServer) error {
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	httpServer, err := NewServer(addr, cacheDir, server)
	if err != nil {
		return err
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-signalCtx.Done()
	// got SIGINT / SIGTERM
	return httpServer.Shutdown(context.TODO())
}
