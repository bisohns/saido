package inspector

import (
	log "github.com/sirupsen/logrus"
)

// DF : Parsing the `df` output for memory monitoring
type DF struct {
	fields
}

func (i *DF) Parse(output string) {
	log.Debug(output)
}

// NewDF : Initialize a new DF instance
func NewDF() *DF {
	return &DF{
		fields: fields{
			Type:    Command,
			Command: `df -a`,
		},
	}

}
