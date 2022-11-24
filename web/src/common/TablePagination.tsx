import { Icon, IconButton } from "@mui/material";
import clsx from "clsx";
// import "./Pagination.css";

/**
 *
 * @param {TablePaginationProps} props
 */
function Pagination(props: any) {
  const { instance, className, classes, ...rest } = props;

  return (
    <div className={clsx("Pagination", className, classes?.root)} {...rest}>
      <span className={clsx("Pagination__info", classes?.info)}>
        {instance.getState().pagination?.pageSize *
          instance.getState().pagination?.pageIndex +
          1}{" "}
        -{" "}
        {instance.getState().pagination?.pageSize *
          (instance.getState().pagination?.pageIndex + 1)}{" "}
        of{" "}
        {instance.options.manualPagination
          ? (instance.options?.pageCount || 0) *
              instance.getState().pagination.pageSize -
            (instance.getState().pagination.pageSize -
              instance.getPrePaginationRowModel().rows.length)
          : instance.getPrePaginationRowModel().rows.length}
      </span>
      <IconButton
        color="inherit"
        size="small"
        onClick={() => instance.setPageIndex(0)}
        disabled={!instance.getCanPreviousPage()}
      >
        <Icon>first_page</Icon>
      </IconButton>
      <IconButton
        color="inherit"
        size="small"
        onClick={() => instance.previousPage()}
        disabled={!instance.getCanPreviousPage()}
      >
        <Icon>navigate_before</Icon>
      </IconButton>
      <div className={clsx("Pagination__page", classes?.page)}>
        <h5 className={clsx("Pagination__pageText", classes?.pageText)}>
          {instance.getState()?.pagination?.pageIndex + 1}
        </h5>
      </div>
      <IconButton
        color="inherit"
        size="small"
        onClick={() => instance.nextPage()}
        disabled={!instance.getCanNextPage()}
      >
        <Icon>navigate_next</Icon>
      </IconButton>
      <IconButton
        color="inherit"
        size="small"
        onClick={() => instance.setPageIndex(instance.getPageCount() - 1)}
        disabled={!instance.getCanNextPage()}
      >
        <Icon>last_page</Icon>
      </IconButton>
    </div>
  );
}

export default Pagination;

/**
 * @typedef {{
 * instance: import("@tanstack/react-table").Table<any>
 * } & import("react").ComponentPropsWithoutRef<'div'>} TablePaginationProps
 */
