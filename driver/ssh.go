package driver

import (
	"fmt"
	"strings"

	"github.com/melbahja/goph"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var port = 22

// SSH : Driver for handling ssh executions
type SSH struct {
	fields
	// User e.g root
	User string
	// Host name/ip e.g 171.23.122.1
	Host string
	// port default is 22
	Port int
	// Path to private key file
	KeyFile string
	// Pass key for key file
	KeyPass string
	// Check known hosts (only disable for tests
	CheckKnownHosts bool
	// set environmental vars for server e.g []string{"DEBUG=1", "FAKE=echo"}
	Vars          []string
	SessionClient *goph.Client
}

func (d *SSH) String() string {
	return fmt.Sprintf("%s (%s)", d.User, d.Host)
}

// set the goph Client
func (d *SSH) Client() (*goph.Client, error) {
	if d.SessionClient == nil {
		var err error
		var client *goph.Client
		var callback ssh.HostKeyCallback
		if d.Port != 0 {
			port = d.Port
		}
		auth, err := goph.Key(d.KeyFile, d.KeyPass)
		if err != nil {
			return nil, err
		}
		if d.CheckKnownHosts {
			callback, err = goph.DefaultKnownHosts()
			if err != nil {
				return nil, err
			}
		} else {
			callback = ssh.InsecureIgnoreHostKey()
		}
		client, err = goph.NewConn(&goph.Config{
			User:     d.User,
			Addr:     d.Host,
			Port:     uint(port),
			Auth:     auth,
			Timeout:  goph.DefaultTimeout,
			Callback: callback,
		})
		d.SessionClient = client
		return client, err
	}
	return d.SessionClient, nil
}

func (d *SSH) ReadFile(path string) (string, error) {
	log.Debugf("Reading remote content %s", path)
	command := fmt.Sprintf(`cat %s`, path)
	return d.RunCommand(command)
}

func (d *SSH) RunCommand(command string) (string, error) {
	// FIXME: provide refresh interface to somehow refresh d.Client if d.SessionClient somehow gets dropped
	log.Debugf("Running remote command %s", command)
	client, err := d.Client()
	if err != nil {
		return ``, err
	}
	defer client.Close()
	if len(d.Vars) != 0 {
		// add env variable to command
		envline := strings.Join(d.Vars, ";")
		command = strings.Join([]string{envline, command}, ";")
	}
	out, err := client.Run(command)
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

func (d *SSH) GetDetails() SystemDetails {
	if d.Info == nil {
		uname, err := d.RunCommand(`uname`)
		// try windows command
		if err != nil {
			windowsName, err := d.RunCommand(`systeminfo | findstr /B /C:"OS Name"`)
			if err == nil {
				if strings.Contains(strings.ToLower(windowsName), "windows") {
					uname = "windows"
				}
			}
		}
		details := &SystemDetails{}
		details.Name = uname
		switch details.Name {
		case "windows":
			details.IsWindows = true
		case "linux":
			details.IsLinux = true
		case "darwin":
			details.IsDarwin = true
		}
		d.Info = details
	}
	return *d.Info
}
