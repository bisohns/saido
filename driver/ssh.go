package driver

import (
	"fmt"
	"github.com/melbahja/goph"
	log "github.com/sirupsen/logrus"
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
}

func (d *SSH) String() string {
	return fmt.Sprintf("%s (%s)", d.User, d.Host)
}

// set the goph Client
func (d *SSH) Client() (*goph.Client, error) {
	port := 22
	if d.Port != 0 {
		port = d.Port
	}
	auth, err := goph.Key(d.PubKeyFile, d.PubKeyPass)
	if err != nil {
		return nil, err
	}
	callback, err := goph.DefaultKnownHosts()
	if err != nil {
		return nil, err
	}
	client, err := goph.NewConn(&goph.Config{
		User:     d.User,
		Addr:     d.Host,
		Port:     uint(port),
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: callback,
	})
	if err != nil {
		return nil, err
	}
	return client, err
}

func (d *SSH) readFile(path string) (string, error) {
	log.Debugf("Reading remote content %s", path)
	command := fmt.Sprintf(`cat %s`, path)
	return d.runCommand(command)
}

func (d *SSH) runCommand(command string) (string, error) {
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

func (d *SSH) getDetails() string {
	return fmt.Sprintf(`SSH - %s`, d.String())
}
