import { getCoreRowModel } from "@tanstack/react-table";
import { useVirtual } from "react-virtual";
import Table from "common/Table";
import useTable from "common/useTable";
import React from "react";

interface ServerDetailServicesTabPanelProcessType {
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

export default function ServerDetailServicesTabPanelProcess(
  props: ServerDetailServicesTabPanelProcessType
) {
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
    header: "First Name",
    accessorKey: "firstname",
  },
  {
    header: "Last Name",
    accessorKey: "lastname",
  },
  {
    header: "Message",
    accessorKey: "message",
  },
  {
    header: "Created on",
    accessorKey: "created_at",
  },
];

const data = Array(1000).fill({
  firstname: "Joseph",
  lastname: "Edache",
  message: "Hello",
  created_at: Date.now(),
});
