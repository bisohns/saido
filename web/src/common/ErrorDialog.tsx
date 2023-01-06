import React from "react";
import {
  Dialog,
  DialogContent,
  DialogContentText,
  DialogTitle,
  DialogActions,
  Button,
} from "@mui/material";

interface ErrorDialogProps {
  title?: string;
  description?: string;
  onRetry?: () => void;
  retryText?: string;
}
/**
 *
 * @param {ErrorDialogProps} props
 */
export function ErrorDialog(props: ErrorDialogProps) {
  const { title, description, onRetry, retryText, ...rest } = props;
  function handleRetry(e: any) {
    e.stopPropagation();
    if (onRetry) {
      onRetry();
    }
  }

  return (
    <Dialog open={true} fullWidth {...rest}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{description}</DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleRetry}>{retryText}</Button>
      </DialogActions>
    </Dialog>
  );
}

ErrorDialog.defaultProps = {
  title: "Something went wrong",
  description:
    "Sorry, something went wrong, Please try again later or contact our support.",
  retryText: "Try Again",
};

export default ErrorDialog;

/**
 * @typedef {import("./ErrorDialogContext").ErrorOptions & import("@material-ui/core").DialogProps} ErrorDialogProps
 */
