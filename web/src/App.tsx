import React, { useEffect } from "react";

import AppThemeProvider from "./AppThemeProvider";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ServerList from "./server/ServerList";
import ServerDetail from "./server/ServerDetail";
import useSocket from "hooks/useSocket";
import ThemeConfig from "ThemeConfig";
import { Box } from "@mui/material";

function App() {
  const { serversGroupedByHost, connectionStatus, setJsonMessage } =
    useSocket();

  useEffect(() => {
    document.body.style.backgroundColor = `${ThemeConfig.palette.common.black}`;
  }, []);

  const router = createBrowserRouter([
    {
      path: "/",
      element: (
        <ServerList
          servers={serversGroupedByHost}
          connectionStatus={connectionStatus}
          setJsonMessage={setJsonMessage}
        />
      ),
    },
    {
      path: "/:host",
      element: (
        <ServerDetail
          servers={serversGroupedByHost}
          connectionStatus={connectionStatus}
          setJsonMessage={setJsonMessage}
        />
      ),
    },
  ]);

  return (
    <AppThemeProvider>
      <div
        style={{
          minHeight: "100vh",
          height: "100%",
        }}
      >
        <RouterProvider router={router} />
        <Box mb={10} />
      </div>
    </AppThemeProvider>
  );
}

export default App;
