package driver

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

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
		d.GetDetails()
	}
	if d.Info.IsLinux || d.Info.IsDarwin {
		cmd = exec.Command("bash", "-c", command)
	} else {
		cmd = exec.Command("cmd", "/C", command)
	}
	cmd.Env = os.Environ()
	if len(d.Vars) != 0 {
		for _, v := range d.Vars {
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
		}
		details.Extra = runtime.GOARCH
		d.Info = details
	}
	return *d.Info
}
