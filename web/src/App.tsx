import React from "react";

import AppThemeProvider from "./AppThemeProvider";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ServerList from "./server/ServerList";
import ServerDetail from "./server/ServerDetail";
import AppHeader from "./AppHeader";

const router = createBrowserRouter([
  {
    path: "/",
    element: <ServerList />,
  },
  {
    path: "/:id",
    element: <ServerDetail />,
  },
]);

function App() {
  return (
    <AppThemeProvider>
      <>
        <AppHeader />
        <RouterProvider router={router} />
      </>
    </AppThemeProvider>
  );
}

export default App;
