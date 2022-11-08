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
  | ProcessData;

export interface ServerResponseType {
  Error: boolean;
  Message: {
    Host: string;
    Name: ServerServiceNameType;
    Platform: 'Windows' | 'Linux' | 'Darwin' | 'MacOS';
    Data: ServerResponseMessageData;
  };
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

export interface DockerData {}

export interface UptimeData {}

export interface ProcessData {}

export interface ServerGroupedByHostResponseType {
  [Host: string]: ServerResponseType[];
}

export interface ServerGroupedByNameResponseType {
  [Name: string]: ServerResponseType[];
}
