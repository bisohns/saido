export type ServerServiceNameType =
  | "disk"
  | "docker"
  | "uptime"
  | "memory"
  | "process"
  | "loadavg"
  | "tcp";

export type ServerResponseMessageData =
  | Array<DiskData>
  | MemoryData
  | Array<DockerData>
  | UptimeData
  | Array<ProcessData>
  | LoadingAvgData
  | TCPData;

export interface ServerResponseType<T = ServerResponseMessageData> {
  Error: boolean;
  Message: {
    Host: string;
    Error?: string;
    Name: ServerServiceNameType;
    Platform: "Windows" | "Linux" | "Darwin" | "MacOS";
    Data: T;
  };
}

export interface ServerResponseByHostType<T = ServerResponseMessageData> {
  [host: string]: Array<{
    Error: boolean;
    Message: {
      Host: string;
      Error?: string;
      Name: ServerServiceNameType;
      Platform: "Windows" | "Linux" | "Darwin" | "MacOS";
      Data: T;
    };
  }>;
}

export interface MemoryData {
  MemFree: number;
  MemTotal: number;
  Cached: number;
  SwapTotal: number;
  SwapFree: number;
}

export interface DiskData {
  Available: number;
  FileSystem: string;
  PercentFull: number;
  Size: number;
  Used: number;
  VolumeName: string;
}

export interface LoadingAvgData {
  Load1M: number;
  Load5M: number;
  Load15M: number;
}

export interface DockerData {
  CPU: 5.52;
  ContainerID: "ab180d7dd324";
  ContainerName: "fastmeet-fast-meet-1";
  Limit: 7851.008;
  MemPercent: 45.46;
  MemUsage: 3568.64;
  Pid: 148;
}

export interface UptimeData {
  Idle: number;
  IdlePercent: number;
  Up: number;
}

export interface TCPData {
  Ports: Record<number, string>;
}

export interface ProcessData {
  Pid: number;
  Memory: number;
  SessionName: string;
  Command: string;
}

export interface ServerGroupedByHostResponseType {
  [Host: string]: ServerResponseType[];
}

export interface ServerGroupedByNameResponseType {
  [Name: string]: {
    data: ServerResponseType[];
    Host: string;
  };
}
