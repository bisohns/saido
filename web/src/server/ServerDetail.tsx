import { Container } from "@mui/material";
import React, { useState } from "react";
import { useParams } from "react-router-dom";
import useWebSocket, { ReadyState } from "react-use-websocket";
import LoadingContent from "../common/LoadingContent";
import {
  ServerGroupedByNameResponseType,
  ServerResponseType,
} from "./ServerType";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import PageHeader from "common/PageHeader";

const wssMetricsBaseURL = `${process.env.REACT_APP_WS_BASE_URL}/metrics`;

export default function ServerDetail() {
  const { host } = useParams<{ host: string }>();
  const [tabIndex, setTabIndex] = React.useState<number>(0);

  const [servers, setServers] = useState<ServerResponseType[]>([]);

  const servicesGroupedByName: ServerGroupedByNameResponseType = servers.reduce(
    (group: any, server) => {
      const { Message } = server;
      const { Name } = Message;
      group[Name] = group[Name] ?? [];
      group[Name].push(server);
      return group;
    },
    {}
  );

  const { sendJsonMessage, readyState } = useWebSocket(wssMetricsBaseURL, {
    onOpen: () => console.log("WebSocket connection opened."),
    onClose: () => console.log("WebSocket connection closed."),
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: WebSocketEventMap["message"]) => {
      const newMessage: ServerResponseType = JSON.parse(event.data);
      if (newMessage.Error) return;
      setServers((prev) => prev.concat(newMessage));
    },
  });

  sendJsonMessage({ FilterBy: host });

  const connectionStatus: string = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  console.log(servicesGroupedByName);

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
        breadcrumbs={[{ name: "Servers", to: "/" }, { name: `${host}` }]}
      ></PageHeader>

      <LoadingContent
        loading={connectionStatus === "Connecting"}
        error={connectionStatus === "Closed"}
      >
        <>
          <Tabs
            value={tabIndex}
            onChange={handleChangeTabIndex}
            aria-label={`${host} Tab`}
            variant="scrollable"
            scrollButtons="auto"
          >
            {Object.keys(servicesGroupedByName)?.map(
              (serverName: string, index: number) => (
                <Tab label={serverName} key={index} />
              )
            )}
          </Tabs>

          {Object.keys(servicesGroupedByName)?.map(
            (serverName: string, index: number) => (
              <div key={index}>{index === tabIndex && <>{serverName}</>}</div>
            )
          )}
        </>
      </LoadingContent>
    </Container>
  );
}
