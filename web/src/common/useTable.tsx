import { useReactTable, TableOptions } from "@tanstack/react-table";

/**
 *
 * @param {import("@tanstack/react-table").TableOptions<any>} options
 * @returns
 */
function useTable(options: TableOptions<any>) {
  return useReactTable({
    debugTable: true,
    ...options,
    data: options.data || [],
  });
}

export default useTable;
