import { Container } from '@mui/material';
import React from 'react';
import { useParams } from 'react-router-dom';
import LoadingContent from '../common/LoadingContent';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import PageHeader from 'common/PageHeader';
import useSocket from 'hooks/useSocket';

export default function ServerDetail() {
  const { host } = useParams<{ host: string }>();

  const [tabIndex, setTabIndex] = React.useState<number>(0);

  const { connectionStatus, sendJsonMessage, servicesGroupedByName } =
    useSocket();
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
            {Object.keys(servicesGroupedByName)?.map(
              (serverName: string, index: number) => (
                <Tab label={serverName} key={index} />
              )
            )}
          </Tabs>
        </>
      </LoadingContent>
    </Container>
  );
}
