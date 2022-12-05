package inspector

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

type TempMetrics struct {
	CPUTemp float64
}

type TempLinux struct {
	Driver            *driver.Driver
	CPUTempFilePath   string
	RawTempDegree     string
	DisplayTempDegree string
	Note              string
	Values            *TempMetrics
}

type TempDarwin struct {
	Driver            *driver.Driver
	CPUTempCommand    string
	RawTempDegree     string
	DisplayTempDegree string
	Note              string
	Values            *TempMetrics
}

type TempWin struct {
	Driver            *driver.Driver
	CPUTempCommand    string
	RawTempDegree     string
	DisplayTempDegree string
	Note              string
	Values            *TempMetrics
}

func (i *TempWin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsWindows {
		panic("Cannot use TempWin outside (windows)")
	}
	i.Driver = driver
}

func (i TempWin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TempWin) Parse(output string) {
	// var values []TempMetrics
	log.Debug("Parsing output string in TempWin inspector")
	lineChar := "\r"
	output = strings.TrimPrefix(output, lineChar)
	output = strings.TrimSuffix(output, lineChar)
	lines := strings.Split(output, "\n")
	if len(lines) < 1 {
		log.Fatalf(`Error Parsing Temperature `)
	}
	farenValue, err := strconv.ParseFloat(lines[1], 64)
	if err != nil {
		log.Fatalf(`Error Parsing Temperature: %s `, err)
	}
	celsValue := (farenValue * 10.0) - 27315.0
	i.Values = &TempMetrics{
		CPUTemp: celsValue,
	}
}

func (i *TempWin) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.CPUTempCommand)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

func (i *TempLinux) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use TempLinux outside (linux)")
	}
	i.Driver = driver
}

func (i TempLinux) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TempLinux) Parse(output string) {
	log.Debug("Parsing output string in Temp inspector")
	output = strings.ReplaceAll(output, " ", "")
	value, err := strconv.ParseFloat(output, 64)
	if err != nil {
		log.Fatalf(`Error Parsing Temperature: %s `, err)
	}
	i.Values = &TempMetrics{
		CPUTemp: value / 1000.0,
	}

}

func (i *TempLinux) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.CPUTempFilePath)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

func (i *TempDarwin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use TempDarwin outside (darwin)")
	}
	i.Driver = driver
}

func (i TempDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TempDarwin) Parse(output string) {
	log.Debug("Parsing output string in Temp inspector")
	tempUnit := "°C"
	output = strings.TrimSuffix(output, tempUnit)
	value, err := strconv.ParseFloat(output, 64)
	if err != nil {
		log.Fatalf(`Error Parsing Temperature: %s `, err)

	}
	i.Values = &TempMetrics{
		CPUTemp: value,
	}

}

func (i *TempDarwin) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.CPUTempCommand)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

func NewTemp(driver *driver.Driver, _ ...string) (Inspector, error) {
	var temp Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use 'temp' command on drivers outside (linux, darwin, windows)")
	}

	// $ cat /sys/class/thermal/thermal_zone*/temp
	// 20000
	// 53000
	// 50000
	// 53000
	// 56000
	// 68000
	// 49000
	// 50000
	// To see what zones the temperatures are referring to use:
	// $ paste <(cat /sys/class/thermal/thermal_zone*/type) <(cat /sys/class/thermal/thermal_zone*/temp) | column -s $'\t' -t | sed 's/\(.\)..$/.\1°C/'
	// INT3400 Thermal  20.0°C
	// SEN1             45.0°C
	// SEN2             51.0°C
	// SEN3             57.0°C
	// SEN4             59.0°C
	// pch_skylake      77.5°C
	// B0D4             50.0°C
	// x86_pkg_temp     51.0°C

	if details.IsLinux {
		temp = &TempLinux{
			CPUTempFilePath:   `/sys/class/thermal/thermal_zone0/temp`,
			RawTempDegree:     `°C`,
			DisplayTempDegree: `°C`,
		}
	} else if details.IsWindows {
		temp = &TempWin{
			CPUTempCommand:    `wmic /namespace:\\root\wmi PATH MSAcpi_ThermalZoneTemperature get CurrentTemperature`,
			RawTempDegree:     `°F`,
			DisplayTempDegree: `°C`,
		}
	} else if details.IsDarwin {
		temp = &TempDarwin{
			CPUTempCommand:    `osx-cpu-temp`,
			RawTempDegree:     `°C`,
			DisplayTempDegree: `°C`,
		}
	}
	temp.SetDriver(driver)
	return temp, nil
}
