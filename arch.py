class Driver:
    def get_details():
        # read /sys/info -> Linux OS
        # read \C\Windows\ -> Windows
        # OS TYPE
        # Flavour (Windows 11, kali linux, Mac Sierra)

class Inspector:
    command = None  # Action
    def __init__(self, driver) 
        self.driver = driver
        self.statistics = None

    def parse(self, output):
        raise NotImplemented

    def _driver_exec(self):
        raise NotImplemented
    
    def execute(self):
        output = None
        if self.driver.is_windows():
            output = self._driver_exec_windows(self.command)
        elif self.driver.is_linux():
            output = self._driver_exec_linux(self.command)
        elif self.driver.is_darwin():
            output = self._driver_exec_darwin(self.command)
        else
            output = self._driver_exec(self.command)
        self.parse(output)

    def set_driver(self, driver):
        self.driver = driver
    

class MemInfoWindows:
    command = "/proc/meminfo"

    def _driver_exec(self):
        return self.driver.read_file

    def _driver_exec(self):
        return self.driver.read_file


class MemInfo(Inspector):
    class __init__(self, driver):
        if driver.is_windows:
            return MemInfoWindows(**)
        if driver.is_windows:
            return MemInfoWindows(**)


class DiskUsage(Inspector):
    command = "df -h"

    def _driver_exec(self):
        return self.driver.run_command


d = SSHDriver(**)
m = MemInfo(d)
m.execute()
print(m.statistics)



