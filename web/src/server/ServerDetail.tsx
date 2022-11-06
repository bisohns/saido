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

ChartJS.register(ArcElement, Tooltip, Legend);

const wssMetricsBaseURL = `${process.env.REACT_APP_WS_BASE_URL}/metrics`;

export default function ServerDetail() {
  const { host } = useParams<{ host: string }>();

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

  const memoryToPieChartData = (serverResponse: ServerResponseType) => {
    const { Message } = serverResponse;
    const { Data, Name } = Message;

    let labels = [];
    let data = [];

    for (const [key, value] of Object.entries(Data)) {
      labels.push(key);
      data.push(value);
    }
    return {
      labels,
      datasets: [
        {
          label: Name,
          data,
          backgroundColor: "#494788",
          borderWidth: 1,
        },
      ],
    };
  };

  return (
    <Container>
      <LoadingContent
        loading={connectionStatus === "Connecting"}
        error={connectionStatus === "Closed"}
      >
        <>
          {Object.keys(serversGroupedByName)?.map(
            (serverName: string, index: number) => (
              <div>
                {serverName === "memory" && (
                  <Pie
                    data={memoryToPieChartData(
                      serversGroupedByName[serverName]?.[0]
                    )}
                  />
                )}
              </div>
            )
          )}
        </>
      </LoadingContent>
    </Container>
  );
}
