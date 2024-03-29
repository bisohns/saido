package driver

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Local : Driver for handling local executions
type Local struct {
	driverBase
	EnvVars []string
}

func (d *Local) ReadFile(path string) (string, error) {
	log.Debugf("Reading content from %s", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return ``, err
	}
	return string(content), nil
}

// RunCommand : For simple commands without shell variables, pipes, e.t.c
// They can be passed directly. For complex commands e.g
// `echo something | awk $var`, turn into a file to be saved
// under ./shell/
func (d *Local) RunCommand(command string) (string, error) {
	// FIXME: If command contains a shell variable $ or glob
	// type pattern, it would not be executed, see
	// https://pkg.go.dev/os/exec for more information
	var cmd *exec.Cmd
	log.Debugf("Running command `%s` ", command)
	if d.Info == nil {
		_, err := d.GetDetails()
		if err != nil {
			return ``, err
		}
	}
	if d.Info.IsLinux || d.Info.IsDarwin {
		cmd = exec.Command("bash", "-c", command)
	} else {
		command = strings.ReplaceAll(command, "\\", "")
		cmd = exec.Command("cmd", "/C", command)
	}
	cmd.Env = os.Environ()
	if len(d.EnvVars) != 0 {
		for _, v := range d.EnvVars {
			cmd.Env = append(cmd.Env, v)
		}
	}
	_ = cmd.Wait()
	out, err := cmd.Output()
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

func (d *Local) GetDetails() (SystemDetails, error) {
	if d.Info == nil {
		details := &SystemDetails{}
		details.Name = strings.Title(runtime.GOOS)
		switch details.Name {
		case "Windows":
			details.IsWindows = true
		case "Linux":
			details.IsLinux = true
		case "Darwin":
			details.IsDarwin = true
		}
		details.Extra = runtime.GOARCH
		d.Info = details
	}
	return *d.Info, nil
}
