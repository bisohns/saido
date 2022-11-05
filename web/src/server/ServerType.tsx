export interface ServerResponse {
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
