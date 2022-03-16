package driver

import (
	"testing"
)

func NewWebForTest() *Web {
	return &Web{
		URL:    "https://duckduckgo.com",
		Method: GET,
		fields: fields{
			PollInterval: 5,
		},
	}
}

func TestWebRunCommand(t *testing.T) {
	d := NewWebForTest()
	output, err := d.RunCommand(`response`)
	if err != nil {
		t.Error(err)
	}
	if output == "" {
		t.Error("Could not parse response time")
	}
}

func TestWebSystemDetails(t *testing.T) {
	d := NewWebForTest()
	details := d.GetDetails()
	if !details.IsWeb {
		t.Errorf("Expected web driver for web test got %s", details.Name)
	}
}
