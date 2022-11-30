import Table from "common/Table";
import useTable from "common/useTable";
import React from "react";
import { getCoreRowModel } from '@tanstack/react-table';
import { ServerResponseType, ServerServiceNameType, TCPData } from "./ServerType";
import { useVirtual } from "react-virtual";

interface ServerDetailServicesTabPanelTCPType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<TCPData>;
}

export default function ServerDetailServicesTabPanelTCP(
  props: ServerDetailServicesTabPanelTCPType
) {
   const {
     serverData: {
       Message: { Data: {Ports } },
     },
   } = props;

   const data = Object.entries(Ports).map(([key, value]) => {
      return {
        Port: key,
        State: value,
      };
   })

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
        variant='default'
        virtualization
        instance={tableInstance}
        virtualizationInstance={rowVirtualizer}
      />
    );
}

const columns = [
  {
    header: 'Port',
    accessorKey: 'Port',
  },
  {
    header: 'State',
    accessorKey: 'State',
  },
];