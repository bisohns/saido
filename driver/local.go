package driver

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Local : Driver for handling local executions
type Local struct {
	fields
	Vars []string
}

func (d *Local) ReadFile(path string) (string, error) {
	log.Debugf("Reading content from %s", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return ``, err
	}
	return string(content), nil
}

func (d *Local) RunCommand(command string) (string, error) {
	// FIXME: If command contains a shell variable $ or glob
	// type pattern, it would not be executed, see
	// https://pkg.go.dev/os/exec for more information
	cmdArgs := strings.Fields(command)
	log.Debugf("Running command `%s` ", command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	if len(d.Vars) != 0 {
		for _, v := range d.Vars {
			cmd.Env = append(cmd.Env, v)
		}
	}
	out, err := cmd.Output()
	fmt.Printf("%s", string(out))
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

func (d *Local) GetDetails() SystemDetails {
	if d.Info == nil {
		details := &SystemDetails{}
		details.Name = runtime.GOOS
		switch details.Name {
		case "windows":
			details.IsWindows = true
		case "linux":
			details.IsLinux = true
		case "darwin":
			details.IsDarwin = true
		default:
			details.IsLinux = true
		}
		details.Extra = runtime.GOARCH
		d.Info = details
	}
	return *d.Info
}
