import React, { useMemo } from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Pie } from "react-chartjs-2";
import {
  ServerResponseType,
  ServerServiceNameType,
  UptimeData,
} from "./ServerType";
import Table from "common/Table";
import useTable from "common/useTable";
import { getCoreRowModel } from "@tanstack/react-table";
import { useVirtual } from "react-virtual";
import { toDaysMinutesSeconds } from "common/utils";

interface ServerDetailServicesTabPanelUptimeType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<UptimeData>;
}

ChartJS.register(ArcElement, Tooltip, Legend);

export default function ServerDetailServicesTabPanelUptime(
  props: ServerDetailServicesTabPanelUptimeType
) {
  const data = [props.serverData?.Message?.Data];

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
    <div>
      <Table
        ref={tableContainerRef}
        variant="default"
        virtualization
        instance={tableInstance}
        virtualizationInstance={rowVirtualizer}
      />
    </div>
  );
}

const columns = [
  {
    header: "Idle",
    accessorFn: (row: UptimeData) => toDaysMinutesSeconds(row.Up),
  },
  {
    header: "IdlePercent",
    accessorFn: (row: UptimeData) => toDaysMinutesSeconds(row.IdlePercent),
  },
  {
    header: "Up",
    accessorFn: (row: UptimeData) => toDaysMinutesSeconds(row.Up),
  },
];
