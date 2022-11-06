export interface ServerResponseType {
  Error: boolean;
  Message: {
    Host: string;
    Name:
      | "disk"
      | "docker"
      | "uptime"
      | "memory"
      | "process"
      | "loadavg"
      | "tcp";
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
