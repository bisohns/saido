import {
  useReactTable,
  getPaginationRowModel,
  TableOptions,
  getCoreRowModel,
} from "@tanstack/react-table";

/**
 *
 * @param {import("@tanstack/react-table").TableOptions<any>} options
 * @returns
 */
function useTable(options: TableOptions<any>) {
  return useReactTable({
    // getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    debugTable: true,
    ...options,
    data: options.data || [],
    // defaultColumn: {
    //   // cell: (info) => <Typography>{info.getValue()}</Typography>,
    //   ...options?.defaultColumn,
    // },
  });
}

export default useTable;
