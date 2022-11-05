import { ThemeProvider, CssBaseline } from "@mui/material";
import theme from "./ThemeConfig";

export function AppThemeProvider(props: { children: JSX.Element }) {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {props.children}
    </ThemeProvider>
  );
}

export default AppThemeProvider;
