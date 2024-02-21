package webdav

import (
	"fmt"
	"net"
	"os/exec"
)

func Mount(addr *net.TCPAddr, location string) error {
	cmd := exec.Command("mount_webdav", "-S", fmt.Sprintf("http://%s", addr.String()), location)
	out, err := cmd.CombinedOutput()
	fmt.Printf("mount: %s\n", out)
	return err
}
