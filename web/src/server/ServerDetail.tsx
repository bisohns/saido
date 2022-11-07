import { Container } from "@mui/material";
import React, { useState } from "react";
import { useParams } from "react-router-dom";
import useWebSocket, { ReadyState } from "react-use-websocket";
import LoadingContent from "../common/LoadingContent";
import {
  ServerGroupedByNameResponseType,
  ServerResponseType,
} from "./ServerType";
import { Pie } from "react-chartjs-2";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import PageHeader from "common/PageHeader";

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

ChartJS.register(ArcElement, Tooltip, Legend);

const wssMetricsBaseURL = `${process.env.REACT_APP_WS_BASE_URL}/metrics`;

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

export default function ServerDetail() {
  const { host } = useParams<{ host: string }>();
  const [tabIndex, setTabIndex] = React.useState<number>(0);

  const [servers, setServers] = useState<ServerResponseType[]>([]);

  const serversGroupedByName: ServerGroupedByNameResponseType = servers.reduce(
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

  console.log(serversGroupedByName);

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
            {Object.keys(serversGroupedByName)?.map(
              (serverName: string, index: number) => (
                <Tab label={serverName} {...a11yProps(index)} key={index} />
              )
            )}
          </Tabs>

          {Object.keys(serversGroupedByName)?.map(
            (serverName: string, index: number) => (
              <div key={index}>{index === tabIndex && <>{serverName}</>}</div>
            )
          )}
        </>
      </LoadingContent>
    </Container>
  );
}
