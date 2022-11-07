import { IconButton, Icon } from "@mui/material";
import clsx from "clsx";
import "./TablePagination.css";

/**
 *
 * @param {TablePaginationProps} props
 */
function TablePagination(props: { instance: any; className: string }) {
  const { instance, className, ...rest } = props;
  const { state, rows } = instance;
  return (
    <div className={clsx("TablePagination", className)} {...rest}>
      <span className="text-paragraph mr-4">
        {state.pageSize * state.pageIndex + 1} -{" "}
        {state.pageSize * (state.pageIndex + 1)} of{" "}
        {instance.manualPagination ? instance.dataCount || 0 : rows.length}
      </span>
      <IconButton
        color="inherit"
        size="small"
        onClick={() => instance.gotoPage(0)}
        disabled={!instance.canPreviousPage}
      >
        <Icon>first_page</Icon>
      </IconButton>
      <IconButton
        color="inherit"
        size="small"
        onClick={instance.previousPage}
        disabled={!instance.canPreviousPage}
      >
        <Icon>navigate_before</Icon>
      </IconButton>
      <div className="rounded w-8 h-8 flex justify-center items-center">
        <h5 className="">{state.pageIndex + 1}</h5>
      </div>
      <IconButton
        color="inherit"
        size="small"
        onClick={instance.nextPage}
        disabled={!instance.canNextPage}
      >
        <Icon>navigate_next</Icon>
      </IconButton>
      <IconButton
        color="inherit"
        size="small"
        onClick={() => instance.gotoPage(instance.pageOptions.length - 1)}
        disabled={!instance.canNextPage}
      >
        <Icon>last_page</Icon>
      </IconButton>
    </div>
  );
}

export default TablePagination;

/**
 * @typedef {{instance: import("react-table").TableInstance} & import("react").ComponentPropsWithoutRef<'div'>} TablePaginationProps
 */
