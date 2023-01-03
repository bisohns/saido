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
  | DockerData
  | UptimeData
  | ProcessData
  | LoadingAvgData;

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

export interface DockerData {}

export interface UptimeData {}

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
