import React from "react";
import { ServerServiceNameType } from "./ServerType";

interface ServerDetailServicesTabPanelLoadAvgType {
  serverName: ServerServiceNameType;
  Error?: string;
  serverData: any;
}

export default function ServerDetailServicesTabPanelCustom(
  props: ServerDetailServicesTabPanelLoadAvgType
) {
  const { serverData } = props;

  if (serverData.Error) {
    return (
      <div>
        <p> {serverData.Message.Error}</p>
      </div>
    );
  }

  return (
    <div>
      <h1>{serverData.Message.Data.Command}</h1>
      <pre className="text-white bg-black p-4">
        {serverData.Message?.Data?.Output}
      </pre>
    </div>
  );
}
