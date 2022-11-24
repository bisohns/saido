import { getCoreRowModel } from "@tanstack/react-table";
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
    getCoreRowModel: getCoreRowModel(),
    initialState: { pagination: { pageSize: 30 } },
    columns,
    data,
  });
  return (
    <div>
      <Table instance={tableInstance} />
    </div>
  );
}

const columns = [
  // {
  //   header: "Name",
  //   footer: (props) => props.column.id,
  //   columns: [
  //     {
  //       accessorKey: "firstname",
  //       cell: (info) => info.getValue(),
  //       footer: (props) => props.column.id
  //     },
  //     {
  //       accessorFn: (row) => row.lastName,
  //       id: "lastname",
  //       cell: (info) => info.getValue(),
  //       header: () => <span>Last Name</span>,
  //       footer: (props) => props.column.id
  //     }
  //   ]
  // },
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

const data = Array(100).fill({
  firstname: "Joseph",
  lastname: "Edache",
  message: "Hello",
  created_at: Date.now(),
});
