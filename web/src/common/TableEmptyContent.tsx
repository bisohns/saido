import { Typography, Icon } from "@mui/material";
import clsx from "clsx";
import "./TableEmptyContent.css";
import FolderOffOutlinedIcon from "@mui/icons-material/FolderOffOutlined";

export function TableEmptyContent(
  props: Partial<{
    title: string;
    className: string;
  }>
) {
  const { title, className, ...rest } = props;
  return (
    <div className={clsx("TableEmptyContent", className)} {...rest}>
      <FolderOffOutlinedIcon className={clsx("TableEmptyContent__icon")} />
      <Typography variant="h6" className={clsx("TableEmptyContent__text")}>
        No data
      </Typography>
    </div>
  );
}

export default TableEmptyContent;
