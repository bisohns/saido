import { ReactNode, useEffect } from "react";
import useDataRef from "hooks/useDataRef";
import LoadingIndicator from "./LoadingIndicator";
import ErrorContent from "common/ErrorContent";
import { Box } from "@mui/material";

interface LoadingContentType {
  size: number;
  error: boolean;
  loading: boolean;
  children: JSX.Element | (() => JSX.Element);
  onReload: () => void;
  onMount: () => void;
  loadingContent: JSX.Element | ((x: JSX.Element) => ReactNode);
  errorContent: JSX.Element | ((x: JSX.Element) => ReactNode);
  className: string;
}
/**
 *
 * @param {LoadingContentProps} props
 */
function LoadingContent(props: Partial<LoadingContentType>): JSX.Element {
  const {
    size,
    error,
    loading,
    children,
    onReload,
    onMount,
    loadingContent,
    errorContent,
    className,
    ...rest
  } = props;

  const dataRef = useDataRef({ onReload, onMount });

  useEffect(() => {
    dataRef.current.onMount?.();
  }, [dataRef]);

  if (!loading && !error) {
    if (children !== undefined) {
      return typeof children === "function" ? children() : children;
    }
  }

  const defaultLoadingContent = <LoadingIndicator size={size} />;

  const defaultErrorContent = <ErrorContent onTryAgain={() => onReload?.()} />;

  return (
    <Box
      display={"flex"}
      justifyContent={"center"}
      alignItems={"center"}
      p={2}
      {...rest}
    >
      {error ? (
        <>
          {errorContent
            ? typeof errorContent === "function"
              ? errorContent(defaultErrorContent)
              : errorContent
            : defaultErrorContent}
        </>
      ) : loadingContent ? (
        typeof loadingContent === "function" ? (
          loadingContent(defaultLoadingContent)
        ) : (
          loadingContent
        )
      ) : (
        defaultLoadingContent
      )}
    </Box>
  );
}

LoadingContent.defaultProps = {
  size: 40,
  children: null,
};

export default LoadingContent;

/**
 * @typedef {{
 * size: string | number,
 * onMount: Function,
 * onReload: Function,
 * error: boolean,
 * loading: boolean,
 * errorContent: React.ReactNode,
 * loadingContent: React.ReactNode,
 * } & React.ComponentPropsWithoutRef<'div'>} LoadingContentProps
 */
