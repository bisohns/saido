import React from "react";

interface ServerDetailServicesTabPanelUptimeType {
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

export default function ServerDetailServicesTabPanelUptime(
  props: ServerDetailServicesTabPanelUptimeType
) {
  return <div>{/* <ServicesTabPanel /> */}</div>;
}
