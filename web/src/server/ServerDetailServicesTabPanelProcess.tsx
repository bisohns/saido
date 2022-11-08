import React from "react";

interface ServerDetailServicesTabPanelProcessType {
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

export default function ServerDetailServicesTabPanelProcess(
  props: ServerDetailServicesTabPanelProcessType
) {
  return <div>{/* <ServicesTabPanel /> */}</div>;
}
