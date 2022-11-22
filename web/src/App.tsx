import React from 'react';

import AppThemeProvider from './AppThemeProvider';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import ServerList from './server/ServerList';
import ServerDetail from './server/ServerDetail';
import AppHeader from './AppHeader';
import useSocket from 'hooks/useSocket';

function App() {
  const { servicesGroupedByName, serversGroupedByHost, updateCount ,connectionStatus,sendJsonMessage} =
    useSocket();

  const [servicesGroupByName, setServicesGroupByName] = React.useState<any>({});

  const [servicesGroupByHost, setServicesGroupByHost] = React.useState<any>({});

  React.useEffect(() => {
    setServicesGroupByName(servicesGroupedByName);
  }, [updateCount]);

  React.useEffect(() => {
    setServicesGroupByHost(serversGroupedByHost);
  }, [updateCount]);

  const router = createBrowserRouter([
    {
      path: '/',
      element: <ServerList serversGroupedByHost={servicesGroupByHost}  connectionStatus={connectionStatus} sendJsonMessage={sendJsonMessage}/>,
    },
    {
      path: '/:host',
      element: <ServerDetail servicesGroupedByName={servicesGroupByName} connectionStatus={connectionStatus} sendJsonMessage={sendJsonMessage} />,
    },
  ]);

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
