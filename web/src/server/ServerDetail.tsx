import React from 'react';
import { Container } from '@mui/material';
import { useParams } from 'react-router-dom';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';

import PageHeader from 'common/PageHeader';
import LoadingContent from '../common/LoadingContent';
import useSocket from 'hooks/useSocket';
import ServerDetailServicesTabPanel from './ServerDetailServicesTabPanel';
import {
  ServerGroupedByNameResponseType,
  ServerResponseType,
  ServerServiceNameType,
} from './ServerType';

export default function ServerDetail({
  servicesGroupedByName,
  connectionStatus,
  sendJsonMessage
}: {
  servicesGroupedByName: ServerGroupedByNameResponseType;
  connectionStatus: string;
  sendJsonMessage:(arg0: any)=>void;
}) {
  const { host } = useParams<{ host: string }>();

  const [tabIndex, setTabIndex] = React.useState<number>(0);
  sendJsonMessage({ FilterBy: host });

  const handleChangeTabIndex = (
    event: React.SyntheticEvent,
    newValue: number
  ) => {
    setTabIndex(newValue);
  };

  return (
    <Container>
      <PageHeader
        title={`${host}`}
        breadcrumbs={[
          { name: 'Servers', to: '/' },
          { name: `${host}` },
        ]}></PageHeader>

      <LoadingContent
        loading={connectionStatus === 'Connecting'}
        error={connectionStatus === 'Closed'}>
        <>
          <Tabs
            value={tabIndex}
            onChange={handleChangeTabIndex}
            aria-label={`${host} Tab`}
            variant='scrollable'
            scrollButtons='auto'>
            {Object.keys(servicesGroupedByName)
              ?.sort()
              ?.map((serverName: string, index: number) => (
                <Tab label={serverName} key={index} />
              ))}
          </Tabs>

          {Object.keys(servicesGroupedByName)
            ?.sort()
            ?.map((serverName: string, index: number) => {
              if (host !== servicesGroupedByName[serverName].Host) return null;

              return (
                <div key={index}>
                  {index === tabIndex && (
                    <ServerDetailServicesTabPanel
                      serverName={serverName as ServerServiceNameType}
                      serverData={
                        servicesGroupedByName[
                          serverName as ServerServiceNameType
                        ]?.data?.at(-1) as ServerResponseType
                      } // get the last object of service
                    />
                  )}
                </div>
              );
            })}
        </>
      </LoadingContent>
    </Container>
  );
}
