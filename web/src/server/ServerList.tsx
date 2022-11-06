import { Card, CardActionArea, Grid, Typography } from "@mui/material";
import { Box, Container } from "@mui/system";
import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { ReactComponent as ServerIcon } from "assets/svg/server.svg";
import LoadingContent from "common/LoadingContent";
import ThemeConfig from "ThemeConfig";
import {
  ServerGroupedByHostResponseType,
  ServerResponseType,
} from "./ServerType";

const wssMetricsBaseURL = `${process.env.REACT_APP_WS_BASE_URL}/metrics`;

export default function ServerList() {
  const navigate = useNavigate();

  const [servers, setServers] = useState<ServerResponseType[]>([]);
  const serversGroupedByHost: ServerGroupedByHostResponseType = servers.reduce(
    (group: any, server) => {
      const { Message } = server;
      const { Host } = Message;
      group[Host] = group[Host] ?? [];
      group[Host].push(server);
      return group;
    },
    {}
  );

  const { readyState } = useWebSocket(wssMetricsBaseURL, {
    onOpen: () => console.log("WebSocket connection opened."),
    onClose: () => console.log("WebSocket connection closed."),
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: WebSocketEventMap["message"]) => {
      const newMessage: ServerResponseType = JSON.parse(event.data);
      setServers((prev) => prev.concat(newMessage));
    },
  });

  const connectionStatus: string = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  return (
    <Container>
      <LoadingContent
        loading={connectionStatus === "Connecting"}
        error={connectionStatus === "Closed"}
      >
        <Grid container spacing={2} my={10}>
          {Object.keys(serversGroupedByHost)?.map(
            (serverHost: string, index: number) => (
              <Grid item xs={6} md={4}>
                <Card key={index}>
                  <CardActionArea onClick={() => navigate(`/${serverHost}`)}>
                    <Box
                      display={"flex"}
                      justifyContent="center"
                      alignItems="center"
                      flexDirection="column"
                    >
                      <ServerIcon width={"100px"} />
                      <Typography
                        textTransform={"capitalize"}
                        mb={2}
                        noWrap
                        fontWeight={600}
                      >
                        <>
                          {serverHost} -{" "}
                          <span
                            style={{ color: ThemeConfig.palette.success.dark }}
                          >
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
