import {
  useTable as useReactTable,
  usePagination,
  useFlexLayout,
  useRowSelect,
  useGlobalFilter,
  UseTableOptions,
  CellValue,
} from "react-table";
import { useSticky } from "react-table-sticky";
import { Typography } from "@mui/material";

const initialData: [] = [];

/**
 * @type {typeof useReactTable}
 */
function useTable(options: UseTableOptions<{}>, ...plugins: []) {
  return useReactTable(
    {
      ...options,
      data: options?.data || initialData,
      defaultColumn: Object.assign(
        {
          Cell: ({ value }: any) => (
            <Typography noWrap textAlign={"left"}>
              {value}
            </Typography>
          ),
        },
        options?.defaultColumn
      ),
    },
    (hook) => {
      hook.prepareRow.push((row: any) => {
        row.normalizedCells = row.cells.reduce(
          (acc: CellValue[], cell: CellValue) => {
            acc[cell.column.id] = cell;
            return acc;
          },
          {}
        );
      });
      hook.visibleColumns.push((visibleColumns: any, { instance }: any) => {
        if (instance.hideRowCounter) {
          return visibleColumns;
        }

        return [
          {
            sticky: "left",
            Header: "#",
            width: 60,
            Cell: (instance) => instance.row.index + 1,
          },
          ...visibleColumns,
        ];
      });
    },
    useFlexLayout,
    useGlobalFilter,
    usePagination,
    useRowSelect,
    useSticky,
    ...plugins
  );
}

export default useTable;
