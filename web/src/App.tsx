import React from "react";

import AppThemeProvider from "./AppThemeProvider";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ServerList from "./server/ServerList";
import ServerDetail from "./server/ServerDetail";
import useSocket from "hooks/useSocket";
import ThemeConfig from "ThemeConfig";

function App() {
  const {
    servicesGroupedByName,
    serversGroupedByHost,
    connectionStatus,
    setJsonMessage,
  } = useSocket();

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
          background: `${ThemeConfig.palette.common.black}`,
          minHeight: "100vh",
        }}
      >
        <RouterProvider router={router} />
      </div>
    </AppThemeProvider>
  );
}

export default App;
