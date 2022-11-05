import React from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

export default function ServerList() {
  const { sendJsonMessage, getWebSocket, readyState } = useWebSocket(
    `ws://localhost:3000/metrics`,
    {
      onOpen: () => console.log("WebSocket connection opened."),
      onClose: () => console.log("WebSocket connection closed."),
      shouldReconnect: (closeEvent) => true,
      onMessage: (event: WebSocketEventMap["message"]) =>
        console.log("new Data", event.data),
    }
  );
  getWebSocket();

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  //   console.log("getWebsocket", getWebSocket()?.OPEN);
  return <div>ServerList</div>;
}
