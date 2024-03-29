package driver

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Request : use specified request methods for web
type Request string

const (
	// POST : HTTP post
	POST Request = "POST"
	// GET : HTTP get
	GET Request = "GET"
)

// Web : Driver for handling web executions
type Web struct {
	driverBase
	// URL e.g https://google.com
	URL string
	// Method POST/GET
	Method Request
	// Payload in case of a POST
	Payload string
}

func (d *Web) String() string {
	return fmt.Sprintf("%s (%s)", d.URL, d.Method)
}

func (d *Web) ReadFile(path string) (string, error) {
	log.Debug("Cannot read file on web driver")
	return ``, errors.New("Cannot read file on web driver")
}

func (d *Web) RunCommand(command string) (string, error) {
	if command == `response` {
		var res *http.Response
		var err error
		start := time.Now()
		if d.Method == POST {
			body := []byte(d.Payload)
			res, err = http.Post(d.URL, "application/json", bytes.NewBuffer(body))
		} else {
			res, err = http.Get(d.URL)
		}
		if err != nil || res.StatusCode < 200 || res.StatusCode > 299 {
			message := fmt.Sprintf("Error %s running request: %d", err, res.StatusCode)
			return ``, errors.New(message)
		}
		elapsed := time.Since(start)
		return strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64), nil
	}
	return ``, errors.New("Cannot read file on web driver")
}

func (d *Web) GetDetails() (SystemDetails, error) {
	if d.Info == nil {
		details := &SystemDetails{
			Name:  "web",
			Extra: d.URL,
			IsWeb: true,
		}
		d.Info = details
	}
	return *d.Info, nil
}
