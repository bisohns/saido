import {
  createTheme,
  responsiveFontSizes,
  ThemeOptions,
} from "@mui/material/styles";

export const theme: any = customizeComponents({
  palette: {
    common: {
      white: "#FFFFFF",
      black: "#000000",
    },
    primary: {
      light: "#9ac4f241",
      main: "#588dc5",
      dark: "#61a3ef8e",
      contrastText: "#FFFFFF",
    },
    secondary: {
      light: "#74CAFF",
      // main: "#1890FF",
      main: "#9cc1e7",
      dark: "#0C53B7",
      contrastText: "#FFFFFF",
    },
    info: {
      light: "#74CAFF",
      main: "#1890FF",
      dark: "#0C53B7",
      contrastText: "#FFFFFF",
    },
    success: {
      light: "#AAF27F",
      main: "#54D62C",
      dark: "#229A16",
      contrastText: "#212B36",
    },
    warning: {
      light: "#FFE16A",
      main: "#FFC107",
      dark: "#B78103",
      contrastText: "#212B36",
    },
    error: {
      light: "#FFA48D",
      main: "#FF4842",
      dark: "#B72136",
      contrastText: "#212B36",
    },
    grey: {
      50: "#fafafa",
      100: "#F9FAFB",
      200: "#F4F6F8",
      300: "#DFE3E8",
      400: "#C4CDD5",
      500: "#919EAB",
      600: "#637381",
      700: "#454F5B",
      800: "#212B36",
      900: "#161C24",
    },
    action: {
      active: "#637381",
      hover: "#919EAB14",
      selected: "#919EAB29",
      disabled: "#919EAB",
      disabledBackground: "#919EAB3D",
      focus: "#919EAB3D",
    },
    text: {
      primary: "#212B36",
      secondary: "#637381",
      disabled: "#919EAB",
    },
    divider: "#919EAB3D",
  },
  breakpoints: {
    values: {
      xs: 0,
      sm: 640,
      md: 768,
      lg: 1024,
      xl: 1280,
    },
  },
  typography: {
    fontFamily: ["Nunito Sans", "sans-serif"].join(),
    fontSize: 12,
    button: {
      textTransform: "none",
    },
  },
  components: {
    MuiTextField: {
      styleOverrides: {
        root: {
          "& .MuiInputBase-root": {
            borderRadius: 8,
          },
        },
      },
    },
    MuiModal: {
      styleOverrides: {
        root: {
          "& .MuiBox-root": {
            borderRadius: 8,
            border: "none",
          },
        },
      },
    },
    MuiPopover: {
      defaultProps: {
        PaperProps: {
          style: {
            borderRadius: 10,
          },
        },
      },
    },
    MuiTabs: {
      defaultProps: {
        variant: "scrollable",
        scrollButtons: "auto",
        allowScrollButtonsMobile: true,
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          padding: 3,
          minWidth: 10,
          marginRight: 30,
        },
      },
      defaultProps: {
        sx: {
          typography: {
            fontSize: 14,
            textTransform: "capitalize",
            fontWeight: 550,
          },
        },
      },
    },
    MuiButton: {
      defaultProps: {
        variant: "contained",
        style: {
          borderRadius: 8,
          border: 0,
        },
      },
    },
  },
});

export default responsiveFontSizes(theme);

/**
 *
 * @param {import("@mui/material").Theme} theme
 */
function customizeComponents(theme: ThemeOptions) {
  return createTheme({
    ...theme,
  });
}
