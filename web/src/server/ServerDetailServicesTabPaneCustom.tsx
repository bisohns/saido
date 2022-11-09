import React from 'react';
import { ServerServiceNameType } from './ServerType';

interface ServerDetailServicesTabPanelLoadAvgType {
  serverName: ServerServiceNameType;
  Error?: string;
  serverData: any;
}

export default function ServerDetailServicesTabPanelCustom(
  props: ServerDetailServicesTabPanelLoadAvgType
) {
  if (props.serverData.Error) {
    return (
      <div>
        <p> {props.serverData.Message.Error}</p>
      </div>
    );
  }

  return (
    <div>
      <pre>{props.serverData.Message?.Data?.Output}</pre>
    </div>
  );
}
