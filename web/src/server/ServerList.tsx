import { Card, CardActionArea, Grid, Typography } from "@mui/material";
import { Box, Container } from "@mui/system";
import React, { useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import { ReactComponent as ServerIcon } from "../assets/svg/server.svg";
import LoadingContent from "../common/LoadingContent";
import ThemeConfig from "../ThemeConfig";
import { ServerResponse } from "./ServerType";

export default function ServerList() {
  const [servers, setServers] = useState<ServerResponse[]>([]);
  const { sendJsonMessage, getWebSocket, readyState } = useWebSocket(
    `ws://localhost:3000/metrics`,
    {
      onOpen: () => console.log("WebSocket connection opened."),
      onClose: () => console.log("WebSocket connection closed."),
      shouldReconnect: (closeEvent) => true,
      onMessage: (event: WebSocketEventMap["message"]) => {
        const newMessage: ServerResponse = JSON.parse(event.data);
        setServers((prev: ServerResponse[]) => {
          if (!newMessage.Error) {
            return prev.concat(newMessage);
          }
        });
      },
    }
  );
  getWebSocket();

  const connectionStatus: string = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  console.log(servers);
  //   console.log("getWebsocket", getWebSocket()?.OPEN);
  return (
    <Container>
      <LoadingContent
        loading={connectionStatus === "Connecting"}
        error={connectionStatus === "Closed"}
      >
        <Grid container spacing={2} my={10}>
          {servers.map((server: any, index: number) => (
            <Grid item xs={6} md={4}>
              <Card key={index}>
                <CardActionArea>
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
                      localhost -{" "}
                      <span style={{ color: ThemeConfig.palette.success.dark }}>
                        linux
                      </span>
                    </Typography>
                  </Box>
                </CardActionArea>
              </Card>
            </Grid>
          ))}
        </Grid>
      </LoadingContent>
    </Container>
  );
}
