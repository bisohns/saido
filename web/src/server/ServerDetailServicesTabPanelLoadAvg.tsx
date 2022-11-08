import React from "react";

interface ServerDetailServicesTabPanelLoadAvgType {
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

export default function ServerDetailServicesTabPanelLoadAvg(
  props: ServerDetailServicesTabPanelLoadAvgType
) {
  return <div>{/* <ServicesTabPanel /> */}</div>;
}
