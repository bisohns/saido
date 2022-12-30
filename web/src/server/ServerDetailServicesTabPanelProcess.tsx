import { getCoreRowModel } from "@tanstack/react-table";
import { useVirtual } from "react-virtual";
import Table from "common/Table";
import useTable from "common/useTable";
import React from "react";
import {
  ProcessData,
  ServerResponseType,
  ServerServiceNameType,
} from "./ServerType";

interface ServerDetailServicesTabPanelProcessType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<ProcessData[]>;
}

export default function ServerDetailServicesTabPanelProcess(
  props: ServerDetailServicesTabPanelProcessType
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
    header: "SessionName",
    accessorKey: "SessionName",
  },
  {
    header: "Memory",
    accessorKey: "Memory",
  },
  {
    header: "Pid",
    accessorKey: "Pid",
  },
  {
    header: "Command",
    accessorKey: "Command",
  },
];
