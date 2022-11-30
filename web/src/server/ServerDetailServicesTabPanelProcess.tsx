import { getCoreRowModel } from "@tanstack/react-table";
import { useVirtual } from "react-virtual";
import Table from "common/Table";
import useTable from "common/useTable";
import React from "react";
import { ProcessData, ServerResponseType } from "./ServerType";

interface ServerDetailServicesTabPanelProcessType {
  serverName:
    | 'disk'
    | 'docker'
    | 'uptime'
    | 'memory'
    | 'process'
    | 'loadavg'
    | 'tcp';
  serverData: ServerResponseType<ProcessData[]>;
}

export default function ServerDetailServicesTabPanelProcess(
  props: ServerDetailServicesTabPanelProcessType
) {
  console.log(props);
  const tableInstance = useTable({
    data:props.serverData.Message.Data,
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
    header: 'Command',
    accessorKey: 'Command',
  },
  {
    header: 'SessionName',
    accessorKey: 'SessionName',
  },
  {
    header: 'Memory',
    accessorKey: 'Memory',
  },
  {
    header: 'Pid',
    accessorKey: 'Pid',
  },
];

