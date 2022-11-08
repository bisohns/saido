export type ServerServiceNameType =
  | "disk"
  | "docker"
  | "uptime"
  | "memory"
  | "process"
  | "loadavg"
  | "tcp";

export interface ServerResponseType {
  Error: boolean;
  Message: {
    Host: string;
    Name: ServerServiceNameType;
    Platform: string;
    Data: Object;
  };
}

export interface ServerGroupedByHostResponseType {
  [Host: string]: ServerResponseType[];
}

export interface ServerGroupedByNameResponseType {
  [Name: string]: ServerResponseType[];
}
