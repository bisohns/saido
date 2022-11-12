import { useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import {
  ServerGroupedByHostResponseType,
  ServerGroupedByNameResponseType,
  ServerResponseType,
} from 'server/ServerType';

const wssMetricsBaseURL = `${process.env.REACT_APP_WS_BASE_URL}/metrics`;

/* 
  This hook is used to connect to the websocket server and send messages to it.
*/
export default function useSocket(options = {}) {
  const [servers, setServers] = useState<ServerResponseType[]>([]);
  const [updateCount, setUpdateCount] = useState<number>(0);

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

  const servicesGroupedByName: ServerGroupedByNameResponseType = servers.reduce(
    (group: any, server: any) => {
      const { Message } = server;
      const { Name } = Message;
      group[Name] = group[Name] ?? [];
      group[Name].push(server);
      return group;
    },
    {}
  );

  //   Uncomment during debugging
  //   console.log('server', servers);

  const { sendJsonMessage, readyState } = useWebSocket(wssMetricsBaseURL, {
    onOpen: () => console.log('WebSocket connection opened.'),
    onClose: () => console.log('WebSocket connection closed.'),
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: WebSocketEventMap['message']) => {
      const newMessage: ServerResponseType = JSON.parse(event.data);
      // if (newMessage.Error) return;
      setUpdateCount((updateCount) => updateCount + 1);
      setServers((prev) => prev.concat(newMessage));
    },
    ...options,
  });

  const connectionStatus: string = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Open',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  return {
    connectionStatus,
    sendJsonMessage,
    servers,
    serversGroupedByHost,
    servicesGroupedByName,
    updateCount,
  };
}
