import { useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import {
  ServerGroupedByHostResponseType,
  ServerGroupedByNameResponseType,
  ServerResponseByHostType,
  ServerResponseType,
} from "server/ServerType";

const wssMetricsBaseURL = `${process.env.REACT_APP_WS_BASE_URL}/metrics`;
const wssMetricsURL = `${
  window.location.protocol === "https:" ? "wss" : "ws"
}://${window.location.host}/metrics`;
/* 
  This hook is used to connect to the websocket server and send messages to it.
*/
export default function useSocket(options = {}) {
  const [servers, setServers] = useState<ServerResponseByHostType>({});
  const [jsonMessage, setJsonMessage] = useState<{ [key: string]: string }>();

  let serversHost: Array<ServerResponseType> =
    servers[jsonMessage?.FilterBy as any] || [];

  const serversGroupedByHost: ServerGroupedByHostResponseType =
    serversHost.reduce((group: any, server) => {
      const { Message } = server;
      const { Host } = Message;
      group[Host] = group[Host] ?? [];
      group[Host].push(server);
      return group;
    }, {});

  const servicesGroupedByName: ServerGroupedByNameResponseType =
    serversHost.reduce((group: any, server: any) => {
      const { Message } = server;
      const { Name, Host } = Message;
      group[Name] = group[Name] ?? { data: [], Host };
      group[Name].data.push(server);
      return group;
    }, {});

  let socketUrl =
    process.env.NODE_ENV === "production" ? wssMetricsURL : wssMetricsBaseURL;

  const { sendJsonMessage, readyState } = useWebSocket(socketUrl, {
    onOpen: () => console.log("WebSocket connection opened."),
    onClose: () => console.log("WebSocket connection closed."),
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: WebSocketEventMap["message"]) => {
      const newMessage: ServerResponseType = JSON.parse(event.data);
      if (jsonMessage) {
        setServers({
          ...servers,
          [jsonMessage?.FilterBy as string]: servers[
            jsonMessage?.FilterBy as string
          ]?.concat(newMessage) || [newMessage],
        });
      }
    },
    ...options,
  });

  useEffect(() => {
    if (jsonMessage) {
      sendJsonMessage(jsonMessage);
    }
  }, [jsonMessage]);

  console.log("servers", servers);

  const connectionStatus: string = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  return {
    connectionStatus,
    setJsonMessage,
    servers,
    serversGroupedByHost,
    servicesGroupedByName,
  };
}
