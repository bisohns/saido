import React from "react";

interface ServerDetailServicesTabPanelDiskType {
  serverName:
    | "disk"
    | "docker"
    | "uptime"
    | "memory"
    | "process"
    | "loadavg"
    | "tcp";
  serverData: Object | [];
}

export default function ServerDetailServicesTabPanelDisk(
  props: ServerDetailServicesTabPanelDiskType
) {
  return <div>{/* <ServicesTabPanel /> */}</div>;
}
