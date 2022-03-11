type Inspector interface {
	getMetrics()
	parse(output string)
	execute()
	_driver_exec() driver.Function
	set_driver(driver)
}

type MemInfoWindows struct{}

type MemInfo struct {
}

func (m *MemInfo) parse()
func (m *MemInfo) execute()
func (m *MemInfo) set_driver()
func (m *MemInfo) _driver_exec()

func NewMemInfo(driver) Inspector {
	switch driver.os {
	case "linux":
		return MemInfoLinux(params)
	}

}

type DiskUsageGeneric struct {
}

type DiskUsageWindows struct {
	*DiskUsageGeneric
}
