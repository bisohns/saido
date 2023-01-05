import { useEffect, useState } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";
import {
  ServerGroupedByHostResponseType,
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

  const serversGroupedByHost: ServerGroupedByHostResponseType = servers;

  const servicesGroupedByName: ServerGroupedByHostResponseType = servers;

  let socketUrl =
    process.env.NODE_ENV === "production" ? wssMetricsURL : wssMetricsBaseURL;

  // console.log("servers", servers);
  const { sendJsonMessage, readyState } = useWebSocket(socketUrl, {
    onOpen: () => console.info("WebSocket connection opened."),
    onClose: () => console.info("WebSocket connection closed."),
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: WebSocketEventMap["message"]) => {
      const newMessage: ServerResponseType = JSON.parse(event.data);

      const newMessageGroupedByHost: ServerGroupedByHostResponseType = [
        newMessage,
      ].reduce((group: any, server) => {
        const { Message } = server;
        const { Host } = Message;
        const serverHost = servers?.[Host] ?? [];

        const serviceIndex = serverHost.findIndex(
          (value) => value.Message.Name === server.Message.Name
        );

        if (serviceIndex > -1) {
          servers[Host][serviceIndex] = server;
        } else {
          group[Host] = [...serverHost, server];
        }

        return group;
      }, {});

      setServers({
        ...servers,
        ...newMessageGroupedByHost,
      });
    },
    ...options,
  });

  useEffect(() => {
    if (jsonMessage) {
      sendJsonMessage(jsonMessage);
    }
  }, [jsonMessage]);

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
