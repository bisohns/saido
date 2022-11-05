import React from "react";
import { Typography } from "@mui/material";
import clsx from "clsx";
import SentimentVeryDissatisfiedIcon from "@mui/icons-material/SentimentVeryDissatisfied";

interface ErrorContentType {
  title: string;
  description: string;
  className?: string;
  onTryAgain?: () => void;
}

/**
 *
 * @param {ErrorContentProps} props
 */

function ErrorContent(props: ErrorContentType): JSX.Element {
  const { title, description, className, onTryAgain, ...rest } = props;

  return (
    <div
      className={clsx(
        "p-1 d-flex text-danger justify-content-center align-items-center flex-column w-100",
        className
      )}
      {...rest}
    >
      <Typography variant="h5" className="font-bold text-center">
        {title}
      </Typography>
      <div>
        <SentimentVeryDissatisfiedIcon fontSize="large" />
      </div>
      <Typography variant="body2" className="text-center mb-4 font-bold w-100">
        {description}
      </Typography>
    </div>
  );
}

ErrorContent.defaultProps = {
  title: "Something went wrong",
  description: "We're quite sorry about this!",
};

export default ErrorContent;

/**
 * @typedef {{
 * onTryAgain: Function
 * } & import("@mui/material").PaperProps} ErrorContentProps
 */
