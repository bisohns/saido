import {
  Cell,
  flexRender,
  Header,
  HeaderGroup,
  Row,
  Table as TanTable,
} from "@tanstack/react-table";
import clsx from "clsx";
import TablePagination from "./TablePagination";
import "./Table.css";
import React, { ReactNode, Ref } from "react";

interface TableProps {
  variant: "default" | "absolute" | "relative";
  instance: TanTable<any>;
  classes?: {
    [P in
      | "root"
      | "row"
      | "table"
      | "empty"
      | "header"
      | "error"
      | "loading"
      | "pagination"
      | "headerRow"
      | "footerRow"
      | "cell"
      | "footerCell"
      | "headerCell"
      | "body"
      | "footer"]: string;
  };
  flexRender: Function;
  virtualizationInstance: any;
  virtualization: boolean;
  onReload: () => void;
  Root: any;
  RootProps: (instance: TanTable<any>, props: TableProps) => any;
  renderRoot: (instance: TanTable<any>, props: TableProps) => ReactNode;
  Table: any;
  TableProps: (instance: TanTable<any>, props: TableProps) => any;
  renderTable: (instance: TanTable<any>, props: TableProps) => ReactNode;
  header: boolean;
  Header: any;
  HeaderProps: (instance: TanTable<any>, props: TableProps) => any;
  renderHeader: (instance: TanTable<any>, props: TableProps) => ReactNode;
  HeaderRow: any;
  HeaderRowProps: (
    headerRow: HeaderGroup<any>,
    instance: TanTable<any>,
    props: TableProps
  ) => any;
  renderHeaderRow: (
    headerRow: HeaderGroup<any>,
    instance: TanTable<any>,
    props: TableProps
  ) => ReactNode;
  HeaderCell: any;
  HeaderCellProps: (
    headerCell: Header<any, any>,
    instance: TanTable<any>,
    props: TableProps
  ) => any;
  renderHeaderCell: (
    headerCell: Header<any, any>,
    instance: TanTable<any>,
    props: TableProps
  ) => ReactNode;
  Body: any;
  BodyProps: (instance: TanTable<any>, props: TableProps) => any;
  renderBody: (instance: TanTable<any>, props: TableProps) => ReactNode;
  BodyRow: any;
  BodyRowProps: (
    bodyRow: Row<any>,
    instance: TanTable<any>,
    props: TableProps
  ) => any;
  renderBodyRow: (
    bodyRow: Row<any>,
    instance: TanTable<any>,
    props: TableProps
  ) => ReactNode;
  BodyCell: any;
  BodyCellProps: (
    bodyCell: Cell<any, any>,
    instance: TanTable<any>,
    props: TableProps
  ) => any;
  renderBodyCell: (
    bodyCell: Cell<any, any>,
    instance: TanTable<any>,
    props: TableProps
  ) => ReactNode;
  footer: boolean;
  Footer: any;
  FooterProps: (instance: TanTable<any>, props: TableProps) => any;
  renderFooter: (instance: TanTable<any>, props: TableProps) => ReactNode;
  FooterRow: any;
  FooterRowProps: (
    footerRow: HeaderGroup<any>,
    instance: TanTable<any>,
    props: TableProps
  ) => any;
  renderFooterRow: (
    footerRow: HeaderGroup<any>,
    instance: TanTable<any>,
    props: TableProps
  ) => ReactNode;
  FooterCell: any;
  FooterCellProps: (
    footerCell: Header<any, any>,
    instance: TanTable<any>,
    props: TableProps
  ) => any;
  renderFooterCell: (
    footerCell: Header<any, any>,
    instance: TanTable<any>,
    props: TableProps
  ) => ReactNode;
  pagination: boolean;
  Pagination: any;
  PaginationProps: (instance: TanTable<any>, props: TableProps) => any;
  renderPagination: (instance: TanTable<any>, props: TableProps) => ReactNode;
  loading: boolean;
  Loading: any;
  LoadingProps: (instance: TanTable<any>, props: TableProps) => any;
  renderLoading: (instance: TanTable<any>, props: TableProps) => ReactNode;
  error: boolean;
  Error: any;
  ErrorProps: (instance: TanTable<any>, props: TableProps) => any;
  renderError: (instance: TanTable<any>, props: TableProps) => ReactNode;
  empty: boolean;
  Empty: any;
  EmptyProps: (instance: TanTable<any>, props: TableProps) => any;
  renderEmpty: (instance: TanTable<any>, props: TableProps) => ReactNode;
}

/**
 *
 * @param {TableProps} props
 */
const Table = React.forwardRef((props: any, ref: Ref<HTMLDivElement>) => {
  return (
    <div ref={ref} className="Table__container">
      {props.renderRoot(props.instance, props)}
    </div>
  );
});

/**
 * @type {TableProps}
 */
Table.defaultProps = {
  header: true,
  footer: false,
  pagination: false,
  virtualization: false,
  flexRender,
  renderRoot,
  renderTable,
  renderHeader,
  renderHeaderRow,
  renderHeaderCell,
  renderBody,
  renderBodyRow,
  renderBodyCell,
  renderFooter,
  renderFooterRow,
  renderFooterCell,
  renderPagination,
  renderLoading,
  renderError,
  renderEmpty,
};

export default Table;

/**
 * @type {TableProps['renderRoot']}
 */
function renderRoot(instance: TanTable<any>, props: TableProps) {
  const {
    classes,
    loading,
    renderLoading,
    error,
    renderError,
    empty = !instance.getPaginationRowModel().rows?.length,
    renderEmpty,
    renderTable,
    pagination,
    renderPagination,
  } = props;

  // const isDefault = props.variant === "default";
  // const isAbsolute = props.variant === "absolute";
  // const isRelative = props.variant === "relative";
  const Root = props.Root || "div";

  return (
    <Root
      {...{
        ...props.RootProps,
        className: clsx(
          "Table",
          // props.RootProps?.className,
          classes?.root
        ),
      }}
    >
      {renderTable?.(instance, props)}
      {loading
        ? renderLoading?.(instance, props)
        : error
        ? renderError?.(instance, props)
        : empty
        ? renderEmpty?.(instance, props)
        : null}
      {pagination && renderPagination?.(instance, props)}
    </Root>
  );
}

/**
 * @type {TableProps['renderTable']}
 */
function renderTable(instance: TanTable<any>, props: TableProps) {
  const {
    classes,
    header,
    footer,
    loading,
    error,
    empty = !instance.getPaginationRowModel().rows?.length,
    renderHeader,
    renderBody,
    renderFooter,
  } = props;

  const isDefault = props.variant === "default";
  const Table = props.Table || (isDefault ? "table" : "div");

  const TableProps = props.TableProps?.(instance, props);

  return (
    <Table
      {...{
        ...TableProps,
        className: clsx("Table__table", TableProps?.className, classes?.table),
        style: {
          width: instance.getTotalSize(),
          ...TableProps?.style,
        },
      }}
    >
      {TableProps?.children || (
        <>
          {header && renderHeader?.(instance, props)}
          {!(loading || error || empty) && (
            <>
              {renderBody?.(instance, props)}
              {footer && renderFooter?.(instance, props)}
            </>
          )}
        </>
      )}
    </Table>
  );
}

/**
 * @type {TableProps['renderHeader']}
 */
function renderHeader(instance: TanTable<any>, props: TableProps) {
  const isDefault = props.variant === "default";
  // const isAbsolute = props.variant === "absolute";
  // const isRelative = props.variant === "relative";
  const Header = props.Header || isDefault ? "thead" : "div";
  const HeaderProps = props.HeaderProps?.(instance, props);

  return (
    <Header
      {...{
        ...HeaderProps,
        className: clsx(
          "Table__table__header",
          HeaderProps?.className,
          props.classes?.header
        ),
      }}
    >
      {HeaderProps?.children ||
        instance
          .getHeaderGroups()
          .map((headerRow) =>
            props.renderHeaderRow(headerRow, instance, props)
          )}
    </Header>
  );
}

/**
 * @type {TableProps['renderHeaderRow']}
 */
function renderHeaderRow(
  headerRow: HeaderGroup<any>,
  instance: TanTable<any>,
  props: TableProps
) {
  const isDefault = props.variant === "default";
  const isAbsolute = props.variant === "absolute";
  const isRelative = props.variant === "relative";
  const HeaderRow = props.HeaderRow || isDefault ? "tr" : "div";
  const HeaderRowProps = props.HeaderRowProps?.(headerRow, instance, props);

  return (
    <HeaderRow
      {...{
        key: headerRow.id,
        ...HeaderRowProps,
        className: clsx(
          "Table__table__header__row",
          props.classes?.headerRow,
          HeaderRowProps?.className,
          isRelative && "Table__table__header__row--relative",
          isAbsolute && "Table__table__header__row--absolute"
        ),
      }}
    >
      {HeaderRowProps?.children ||
        headerRow.headers.map((headerCell) =>
          props.renderHeaderCell(headerCell, instance, props)
        )}
    </HeaderRow>
  );
}

/**
 * @type {TableProps['renderHeaderCell']}
 */
function renderHeaderCell(
  headerCell: Header<any, any>,
  instance: TanTable<any>,
  props: TableProps
) {
  const isDefault = props.variant === "default";
  const isAbsolute = props.variant === "absolute";
  const isRelative = props.variant === "relative";
  const HeaderCell = props.HeaderCell || isDefault ? "th" : "div";
  const HeaderCellProps = props.HeaderCellProps?.(headerCell, instance, props);

  return (
    <HeaderCell
      {...{
        key: headerCell.id,
        colSpan: headerCell.colSpan,
        ...HeaderCellProps,
        className: clsx(
          "Table__table__header__row__cell",
          props.classes?.headerCell,
          HeaderCellProps?.className,
          isRelative && "Table__table__header__row__cell--relative",
          isAbsolute && "Table__table__header__row__cell--absolute"
        ),
        style: {
          width: headerCell.getSize(),
          ...(isAbsolute ? { left: headerCell.getStart() } : {}),
          ...HeaderCellProps?.style,
        },
      }}
    >
      {HeaderCellProps?.children || (
        <>
          {headerCell.isPlaceholder
            ? null
            : props.flexRender(
                headerCell.column.columnDef.header,
                headerCell.getContext()
              )}
          <div
            {...{
              onMouseDown: headerCell.getResizeHandler(),
              onTouchStart: headerCell.getResizeHandler(),
              className: `Table__resizer ${
                headerCell.column.getIsResizing() ? "Table__resizing" : ""
              }`,
              style: {
                transform:
                  instance.options.columnResizeMode === "onEnd" &&
                  headerCell.column.getIsResizing()
                    ? `translateX(${
                        instance.getState().columnSizingInfo.deltaOffset
                      }px)`
                    : "",
              },
            }}
          />
        </>
      )}
    </HeaderCell>
  );
}

/**
 * @type {TableProps['renderBody']}
 */
function renderBody(instance: TanTable<any>, props: TableProps) {
  const isDefault = props.variant === "default";
  // const isAbsolute = props.variant === "absolute";
  // const isRelative = props.variant === "relative";
  const Body = props.Body || isDefault ? "tbody" : "div";
  const BodyProps = props.BodyProps?.(instance, props);
  const { rows } = instance.getRowModel();
  console.log("rows.length", rows);
  const { virtualItems: virtualRows, totalSize } = props.virtualizationInstance;
  const paddingTop = virtualRows.length > 0 ? virtualRows?.[0]?.start || 0 : 0;
  const paddingBottom =
    virtualRows.length > 0
      ? totalSize - (virtualRows?.[virtualRows.length - 1]?.end || 0)
      : 0;

  return (
    <Body
      {...{
        ...BodyProps,
        className: clsx(
          "Table__table__body",
          BodyProps?.className,
          props.classes?.body
        ),
      }}
    >
      <>
        {props.virtualization ? (
          <>
            {paddingTop > 0 && (
              <tr>
                <td style={{ height: `${paddingTop}px` }} />
              </tr>
            )}
            {virtualRows.map((virtualRow: any) =>
              props.renderBodyRow(rows[virtualRow.index], instance, props)
            )}
            {paddingBottom > 0 && (
              <tr>
                <td style={{ height: `${paddingBottom}px` }} />
              </tr>
            )}
          </>
        ) : (
          BodyProps?.children ||
          instance
            .getPaginationRowModel()
            .rows.map((bodyRow) =>
              props.renderBodyRow(bodyRow, instance, props)
            )
        )}
      </>
    </Body>
  );
}

/**
 * @type {TableProps['renderBodyRow']}
 */
function renderBodyRow(
  bodyRow: Row<any>,
  instance: TanTable<any>,
  props: TableProps
) {
  const isAbsolute = props.variant === "absolute";
  const isRelative = props.variant === "relative";
  const Row = props.BodyRow || props.variant === "default" ? "tr" : "div";
  const BodyRowProps = props.BodyRowProps?.(bodyRow, instance, props);

  return (
    <Row
      {...{
        key: bodyRow.id,
        ...BodyRowProps,
        className: clsx(
          "Table__table__body__row",
          BodyRowProps?.className,
          props.classes?.row,
          isRelative && "Table__table__body__row--relative",
          isAbsolute && "Table__table__body__row--absolute"
        ),
      }}
    >
      {BodyRowProps?.children ||
        bodyRow
          .getVisibleCells()
          .map((bodyCell) => props.renderBodyCell(bodyCell, instance, props))}
    </Row>
  );
}

/**
 * @type {TableProps['renderBodyCell']}
 */
function renderBodyCell(
  bodyCell: Cell<any, any>,
  instance: TanTable<any>,
  props: TableProps
) {
  const isAbsolute = props.variant === "absolute";
  const isRelative = props.variant === "relative";
  const Cell = props.BodyCell || props.variant === "default" ? "td" : "div";
  const BodyCellProps = props.BodyCellProps?.(bodyCell, instance, props);

  return (
    <Cell
      {...{
        key: bodyCell.id,
        ...BodyCellProps,
        style: {
          width: bodyCell.column.getSize(),
          ...(isAbsolute ? { left: bodyCell.column.getStart() } : {}),
          ...BodyCellProps?.style,
        },
        className: clsx(
          "Table__table__body__row__cell",
          BodyCellProps?.className,
          props.classes?.cell,
          isRelative && "Table__table__body__row__cell--relative",
          isAbsolute && "Table__table__body__row__cell--absolute"
        ),
      }}
    >
      {BodyCellProps?.children ||
        props.flexRender(bodyCell.column.columnDef.cell, bodyCell.getContext())}
    </Cell>
  );
}

/**
 * @type {TableProps['renderFooter']}
 */
function renderFooter(instance: TanTable<any>, props: TableProps) {
  const isDefault = props.variant === "default";
  // const isAbsolute = props.variant === "absolute";
  // const isRelative = props.variant === "relative";
  const Footer = props.Footer || isDefault ? "tfoot" : "div";
  const FooterProps = props.FooterProps?.(instance, props);

  return (
    <Footer
      {...{
        ...FooterProps,
        className: clsx(
          "Table__table__footer",
          FooterProps?.className,
          props.classes?.footer
        ),
      }}
    >
      {FooterProps?.children ||
        instance
          .getFooterGroups()
          .map((footerRow) =>
            props.renderFooterRow(footerRow, instance, props)
          )}
    </Footer>
  );
}

/**
 * @type {TableProps['renderFooterRow']}
 */
function renderFooterRow(
  footerRow: HeaderGroup<any>,
  instance: TanTable<any>,
  props: TableProps
) {
  const isDefault = props.variant === "default";
  const isAbsolute = props.variant === "absolute";
  const isRelative = props.variant === "relative";
  const FooterRow = props.FooterRow || isDefault ? "tfoot" : "div";
  const FooterRowProps = props.FooterRowProps?.(footerRow, instance, props);

  return (
    <FooterRow
      {...{
        key: footerRow.id,
        ...FooterRowProps,
        className: clsx(
          "Table__table__footer__row",
          FooterRowProps?.className,
          props.classes?.footerRow,
          isRelative && "Table__table__footer__row--relative",
          isAbsolute && "Table__table__footer__row--absolute"
        ),
      }}
    >
      {FooterRowProps?.children ||
        footerRow.headers.map((footerCell) =>
          props.renderFooterCell(footerCell, instance, props)
        )}
    </FooterRow>
  );
}

/**
 * @type {TableProps['renderFooterCell']}
 */
function renderFooterCell(
  footerCell: Header<any, any>,
  instance: TanTable<any>,
  props: TableProps
) {
  const isDefault = props.variant === "default";
  const isAbsolute = props.variant === "absolute";
  const isRelative = props.variant === "relative";
  const FooterCell = props.FooterCell || isDefault ? "th" : "div";
  const FooterCellProps = props.FooterCellProps?.(footerCell, instance, props);

  return (
    <FooterCell
      {...{
        key: footerCell.id,
        colSpan: footerCell.colSpan,
        ...FooterCellProps,
        className: clsx(
          "Table__table__footer__row__cell",
          ...FooterCellProps?.className,
          props.classes?.footerCell,
          isRelative && "Table__table__footer__row__cell--relative",
          isAbsolute && "Table__table__footer__row__cell--absolute"
        ),
        style: {
          width: footerCell.getSize(),
          ...(isAbsolute ? { left: footerCell.getStart() } : {}),
          ...FooterCellProps?.style,
        },
      }}
    >
      {FooterCellProps?.children || (
        <>
          {footerCell.isPlaceholder
            ? null
            : props.flexRender(
                footerCell.column.columnDef.footer,
                footerCell.getContext()
              )}
        </>
      )}
    </FooterCell>
  );
}

/**
 * @type {TableProps['renderPagination']}
 */
function renderPagination(instance: TanTable<any>, props: TableProps) {
  const Pagination = props.Pagination || "div";
  const PaginationProps = props.PaginationProps?.(instance, props);

  return (
    <Pagination
      {...{
        ...PaginationProps,
        className: clsx(
          "Table__pagination",
          PaginationProps?.className,
          props.classes?.pagination
        ),
      }}
    >
      {PaginationProps?.children || <TablePagination instance={instance} />}
    </Pagination>
  );
}

/**
 * @type {TableProps['renderLoading']}
 */
function renderLoading(instance: TanTable<any>, props: TableProps) {
  const Loading = props.Loading || "div";
  const LoadingProps = props.LoadingProps?.(instance, props);

  return (
    <Loading
      {...{
        ...LoadingProps,
        className: clsx(
          "Table__loading",
          LoadingProps?.className,
          props.classes?.loading
        ),
      }}
    >
      {LoadingProps?.children || "Loading..."}
    </Loading>
  );
}

/**
 * @type {TableProps['renderError']}
 */
function renderError(instance: TanTable<any>, props: TableProps) {
  const Error = props.Error || "div";
  const ErrorProps = props.ErrorProps?.(instance, props);

  return (
    <Error
      {...{
        ...ErrorProps,
        className: clsx(
          "Table__error",
          ErrorProps?.className,
          props.classes?.error
        ),
      }}
    >
      {ErrorProps?.children || (
        <>
          Oops!... Something went wrong{" "}
          <button onClick={props.onReload}>Reload</button>
        </>
      )}
    </Error>
  );
}

/**
 * @type {TableProps['renderEmpty']}
 */
function renderEmpty(instance: TanTable<any>, props: TableProps) {
  const Empty = props.Empty || "div";
  const EmptyProps = props.EmptyProps?.(instance, props);

  return (
    <Empty
      {...{
        ...EmptyProps,
        className: clsx(
          "Table__empty",
          EmptyProps?.className,
          props.classes?.empty
        ),
      }}
    >
      {EmptyProps?.children || "No Data"}
    </Empty>
  );
}

/**
 * @typedef {import("@tanstack/react-table").Table<any>} Table
 */

/**
 * @typedef {import("react").ReactNode} ReactNode
 */

/**
 * @typedef  {{
 * variant: "default" | "absolute" | "relative"
 * instance: Table<any>
 * classes: {[P in 'root' | 'table' | 'header' | 'headerRow' | 'headerCell' | 'body' | 'footer']: string};
 * flexRender: Function;
 * onReload: Function ;
 * Root: any;
 * RootProps: (instance: TanTable<any>, props: TableProps) => any;
 * renderRoot: (instance: TanTable<any>, props: TableProps) => ReactNode;
 * Table: any;
 * TableProps: (instance: TanTable<any>, props: TableProps) => any;
 * renderTable: (instance: TanTable<any>, props: TableProps) => ReactNode;
 * header: boolean;
 * Header: any;
 * HeaderProps: (instance: TanTable<any>, props: TableProps) => any;
 * renderHeader: (instance: TanTable<any>, props: TableProps) => ReactNode;
 * HeaderRow: any;
 * HeaderRowProps: (headerRow: HeaderGroup<any>, instance: TanTable<any>, props: TableProps) => any;
 * renderHeaderRow: (headerRow: HeaderGroup<any>, instance: TanTable<any>, props: TableProps) => ReactNode;
 * HeaderCell: any;
 * HeaderCellProps: (headerCell: Header<any, any>, instance: TanTable<any>, props: TableProps) => any;
 * renderHeaderCell: (headerCell: Header<any, any>, instance: TanTable<any>, props: TableProps) => ReactNode;
 * Body: any;
 * BodyProps: (instance: TanTable<any>, props: TableProps) => any;
 * renderBody: (instance: TanTable<any>, props: TableProps) => ReactNode;
 * BodyRow: any;
 * BodyRowProps: (bodyRow: Row<any>, instance: TanTable<any>, props: TableProps) => any;
 * renderBodyRow: (bodyRow: Row<any>, instance: TanTable<any>, props: TableProps) => ReactNode;
 * BodyCell: any;
 * BodyCellProps: (bodyCell: Cell<any>, instance: TanTable<any>, props: TableProps) => any;
 * renderBodyCell: (bodyCell: Cell<any>, instance: TanTable<any>, props: TableProps) => ReactNode;
 * footer: boolean;
 * Footer: any;
 * FooterProps: (instance: TanTable<any>, props: TableProps) => any;
 * renderFooter: (instance: TanTable<any>, props: TableProps) => ReactNode;
 * FooterRow: any;
 * FooterRowProps: (footerRow: HeaderGroup<any>, instance: TanTable<any>, props: TableProps) => any;
 * renderFooterRow: (footerRow: HeaderGroup<any>, instance: TanTable<any>, props: TableProps) => ReactNode;
 * FooterCell: any;
 * FooterCellProps: (footerCell: Header<any, any>, instance: TanTable<any>, props: TableProps) => any;
 * renderFooterCell: (footerCell: Header<any, any>, instance: TanTable<any>, props: TableProps) => ReactNode;
 * pagination: boolean;
 * Pagination: any;
 * PaginationProps: (instance: TanTable<any>, props: TableProps)  => any;
 * renderPagination: (instance: TanTable<any>, props: TableProps)  => ReactNode;
 * loading: boolean;
 * Loading: any;
 * LoadingProps: (instance: TanTable<any>, props: TableProps)  => any;
 * renderLoading: (instance: TanTable<any>, props: TableProps)  => ReactNode;
 * error: boolean;
 * Error: any;
 * ErrorProps: (instance: TanTable<any>, props: TableProps)  => any;
 * renderError: (instance: TanTable<any>, props: TableProps)  => ReactNode;
 * empty: boolean;
 * Empty: any;
 * EmptyProps: (instance: TanTable<any>, props: TableProps)  => any;
 * renderEmpty: (instance: TanTable<any>, props: TableProps)  => ReactNode;
 * }} TableProps
 */
