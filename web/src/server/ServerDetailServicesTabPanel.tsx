import React from "react";
import { ServerNameEnum } from "./ServerConstant";

interface ServerDetailServicesTabPanelType {
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

export default function ServerDetailServicesTabPanel(
  props: ServerDetailServicesTabPanelType
) {
  //   const ServicesTabPanel: React.ElementType = {
  //     [ServerNameEnum.DISK as "disk"]: <></>,
  //     [ServerNameEnum.DOCKER]: <></>,
  //     [ServerNameEnum.LOAD_AVG]: <></>,
  //     [ServerNameEnum.MEMORY]: <></>,
  //     [ServerNameEnum.PROCESS]: <></>,
  //     [ServerNameEnum.TCP]: <></>,
  //     [ServerNameEnum.UPTIME]: <></>,
  //   }[ServerNameEnum];

  return <div>{/* <ServicesTabPanel /> */}</div>;
}
