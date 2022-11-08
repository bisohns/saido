import React from "react";

interface ServerDetailServicesTabPanelMemoryType {
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

export default function ServerDetailServicesTabPanelMemory(
  props: ServerDetailServicesTabPanelMemoryType
) {
  return <div>{/* <ServicesTabPanel /> */}</div>;
}
