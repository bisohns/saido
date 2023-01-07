import { CircularProgress } from "@mui/material";
import React from "react";

/**
 *
 * @param {import("@mui/material").CircularProgressProps} props
 */
function LoadingIndicator(props: { [rest: string]: any }) {
  return <CircularProgress {...props}></CircularProgress>;
}

export default LoadingIndicator;
