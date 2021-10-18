package driver

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

// Local : Driver for handling local executions
type Local struct{}

func (d *Local) readFile(path string) (string, error) {
	log.Debugf("Reading content from %s", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return ``, err
	}
	return string(content), nil
}

func (d *Local) runCommand(command string) (string, error) {
	cmdArgs := strings.Fields(command)
	log.Debugf("Running command `%s` ", command)
	out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Output()
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

func (d *Local) getDetails() string {
	return fmt.Sprintf(`Local - %s`, runtime.GOOS)
}
