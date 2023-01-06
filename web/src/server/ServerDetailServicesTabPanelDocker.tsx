import React from "react";
import {
  DockerData,
  ServerResponseType,
  ServerServiceNameType,
} from "./ServerType";
import { getCoreRowModel } from "@tanstack/react-table";
import { useVirtual } from "react-virtual";
import Table from "common/Table";
import useTable from "common/useTable";
interface ServerDetailServicesTabPanelDockerType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<Array<DockerData>>;
}

export default function ServerDetailServicesTabPanelDocker(
  props: ServerDetailServicesTabPanelDockerType
) {
  const {
    serverData: {
      Message: { Data: data },
    },
  } = props;
  const tableInstance = useTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const tableContainerRef = React.useRef<HTMLDivElement>(null);

  const { rows } = tableInstance.getRowModel();

  const rowVirtualizer = useVirtual({
    parentRef: tableContainerRef,
    size: rows.length,
    overscan: 10,
  });
  return (
    <Table
      ref={tableContainerRef}
      variant="default"
      virtualization
      instance={tableInstance}
      virtualizationInstance={rowVirtualizer}
    />
  );
}

const columns = [
  {
    header: "CPU",
    accessorKey: "CPU",
  },
  {
    header: "ContainerID",
    accessorKey: "ContainerID",
  },
  {
    header: "ContainerName",
    accessorKey: "ContainerName",
  },
  {
    header: "Limit",
    accessorKey: "Limit",
  },
  {
    header: "MemPercent",
    accessorKey: "MemPercent",
  },
  {
    header: "MemUsage",
    accessorKey: "MemUsage",
  },
  {
    header: "Pid",
    accessorKey: "Pid",
  },
];
