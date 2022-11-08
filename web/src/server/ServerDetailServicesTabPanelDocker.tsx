import React from "react";

interface ServerDetailServicesTabPanelDockerType {
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

export default function ServerDetailServicesTabPanelDocker(
  props: ServerDetailServicesTabPanelDockerType
) {
  return <div>{/* <ServicesTabPanel /> */}</div>;
}
