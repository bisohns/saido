package inspector

import (
	log "github.com/sirupsen/logrus"
)

// CustomMetrics : Metrics used by Custom
type CustomMetrics struct {
	Output string
}

// Custom : Parsing the custom command output for disk monitoring
type Custom struct {
	fields
	Values CustomMetrics
}

// Parse : run custom parsing on output of the command
func (i *Custom) Parse(output string) {
	log.Debug("Parsing ouput string in Custom inspector")
	i.Values = i.createMetric(output)
}

func (i Custom) createMetric(output string) CustomMetrics {
	return CustomMetrics{
		Output: output,
	}
}

// NewCustom : Initialize a new Custom instance
func NewCustom(custom string) *Custom {
	return &Custom{
		fields: fields{
			Type:    Command,
			Command: custom,
		},
	}

}
