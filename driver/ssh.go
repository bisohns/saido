package driver

import (
	"fmt"

	"github.com/melbahja/goph"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SSH : Driver for handling ssh executions
type SSH struct {
	fields
	// User e.g root
	User string
	// Host name/ip e.g 171.23.122.1
	Host string
	// port default is 22
	Port int
	// Path to public key file
	PubKeyFile string
	// Pass key for public key file
	PubKeyPass string
	// Check known hosts (only disable for tests
	CheckKnownHosts bool
}

func (d *SSH) String() string {
	return fmt.Sprintf("%s (%s)", d.User, d.Host)
}

// set the goph Client
func (d *SSH) Client() (*goph.Client, error) {
	var err error
	var client *goph.Client
	var callback ssh.HostKeyCallback
	port := 22
	if d.Port != 0 {
		port = d.Port
	}
	auth, err := goph.Key(d.PubKeyFile, d.PubKeyPass)
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
	return client, err
}

func (d *SSH) ReadFile(path string) (string, error) {
	log.Debugf("Reading remote content %s", path)
	command := fmt.Sprintf(`cat %s`, path)
	return d.RunCommand(command)
}

func (d *SSH) RunCommand(command string) (string, error) {
	// FIXME: Do we retain client across all command runs?
	log.Debugf("Running remote command %s", command)
	client, err := d.Client()
	if err != nil {
		return ``, err
	}
	defer client.Close()
	out, err := client.Run(command)
	if err != nil {
		return ``, err
	}
	return string(out), nil
}

func (d *SSH) GetDetails() string {
	return fmt.Sprintf(`SSH - %s`, d.String())
}
