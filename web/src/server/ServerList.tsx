import { Card, CardActionArea, Grid, Typography } from '@mui/material';
import { Box, Container } from '@mui/system';
import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ReactComponent as ServerIcon } from 'assets/svg/server.svg';
import LoadingContent from 'common/LoadingContent';
import ThemeConfig from 'ThemeConfig';
import useSocket from 'hooks/useSocket';

export default function ServerList() {
  const navigate = useNavigate();

  const { connectionStatus, sendJsonMessage, serversGroupedByHost } =
    useSocket();

  useEffect(() => {
    sendJsonMessage({ FilterBy: '' });
  }, []);

  return (
    <Container>
      <LoadingContent
        loading={connectionStatus === 'Connecting'}
        error={connectionStatus === 'Closed'}>
        <Grid container spacing={2} my={10}>
          {Object.keys(serversGroupedByHost)?.map(
            (serverHost: string, index: number) => (
              <Grid item xs={6} md={4}>
                <Card key={index}>
                  <CardActionArea onClick={() => navigate(`/${serverHost}`)}>
                    <Box
                      display={'flex'}
                      justifyContent='center'
                      alignItems='center'
                      flexDirection='column'>
                      <ServerIcon width={'100px'} />
                      <Typography
                        textTransform={'capitalize'}
                        mb={2}
                        noWrap
                        fontWeight={600}>
                        <>
                          {serverHost} -{' '}
                          <span
                            style={{ color: ThemeConfig.palette.success.dark }}>
                            {
                              serversGroupedByHost[serverHost]?.[0]?.Message
                                ?.Platform
                            }
                          </span>
                        </>
                      </Typography>
                    </Box>
                  </CardActionArea>
                </Card>
              </Grid>
            )
          )}
        </Grid>
      </LoadingContent>
    </Container>
  );
}
