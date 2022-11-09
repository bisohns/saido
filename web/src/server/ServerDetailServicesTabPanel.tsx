import React, { useMemo } from "react";

import { ServerNameEnum } from "./ServerConstant";
import ServerDetailServicesTabPanelCustom from "./ServerDetailServicesTabPaneCustom";
import ServerDetailServicesTabPanelDisk from "./ServerDetailServicesTabPanelDisk";
import ServerDetailServicesTabPanelDocker from "./ServerDetailServicesTabPanelDocker";
import ServerDetailServicesTabPanelLoadAvg from "./ServerDetailServicesTabPanelLoadAvg";
import ServerDetailServicesTabPanelMemory from "./ServerDetailServicesTabPanelMemory";
import ServerDetailServicesTabPanelProcess from "./ServerDetailServicesTabPanelProcess";
import ServerDetailServicesTabPanelTCP from "./ServerDetailServicesTabPanelTCP";
import ServerDetailServicesTabPanelUptime from "./ServerDetailServicesTabPanelUptime";
import { DiskData, LoadingAvgData, MemoryData, ServerResponseType, ServerServiceNameType } from "./ServerType";

interface ServerDetailServicesTabPanelType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType;
}

interface servicesTabPanelType {
  title: ServerServiceNameType;
  content: React.ReactElement;
}

export default function ServerDetailServicesTabPanel(
  props: ServerDetailServicesTabPanelType
) {
  const { serverData, serverName } = props;

  const servicesTabPanel: servicesTabPanelType[] = useMemo(
    () => [
      {
        title: ServerNameEnum.DISK as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelDisk
            serverName={serverName}
            serverData={serverData as ServerResponseType<Array<DiskData>>}
          />
        ),
      },
      {
        title: ServerNameEnum.DOCKER as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelDocker
            serverName={serverName}
            serverData={serverData}
          />
        ),
      },
      {
        title: ServerNameEnum.LOAD_AVG as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelLoadAvg
            serverName={serverName}
            serverData={serverData as ServerResponseType<LoadingAvgData>}
          />
        ),
      },
      {
        title: ServerNameEnum.MEMORY as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelMemory
            serverName={serverName}
            serverData={serverData as ServerResponseType<MemoryData>}
          />
        ),
      },
      {
        title: ServerNameEnum.PROCESS as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelProcess
            serverName={serverName}
            serverData={serverData}
          />
        ),
      },
      {
        title: ServerNameEnum.TCP as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelTCP
            serverName={serverName}
            serverData={serverData}
          />
        ),
      },
      {
        title: ServerNameEnum.UPTIME as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelUptime
            serverName={serverName}
            serverData={serverData}
          />
        ),
      },
      {
        title: ServerNameEnum.CUSTOM as ServerServiceNameType,
        content: (
          <ServerDetailServicesTabPanelCustom
            serverName={serverName}
            serverData={serverData}
          />
        ),
      },
    ],
    [serverName, serverData]
  );

  const activeServicesTabPanel = servicesTabPanel?.find(
    (service: servicesTabPanelType) => service.title === serverName
  );

  return <div>{activeServicesTabPanel?.content}</div>;
}
