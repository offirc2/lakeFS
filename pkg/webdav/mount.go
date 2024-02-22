package webdav

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	MountError = errors.New("mount command failed")
)

func execMountCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmdText := fmt.Sprintf("%s %s", name, strings.Join(args, " "))
		return fmt.Errorf("%w: \"%s\":\n%s\n%s", MountError, cmdText, out, err)
	}
	return nil
}
