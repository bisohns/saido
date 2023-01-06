import styled from "@emotion/styled";
import React from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import {
  MemoryData,
  ServerResponseType,
  ServerServiceNameType,
} from "./ServerType";

interface ServerDetailServicesTabPanelMemoryType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<MemoryData>;
}

const Div = styled.div`
  margin-top: 2rem;
`;

export default function ServerDetailServicesTabPanelMemory(
  props: ServerDetailServicesTabPanelMemoryType
) {
  const {
    serverData: {
      Message: { Data },
    },
  } = props;
  return (
    <Div>
      <ResponsiveContainer width="100%" height={500}>
        <BarChart data={[Data]}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="FileSystem" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey="MemFree" fill="#003f5c" />
          <Bar dataKey="MemTotal" fill="#58508d" />
          <Bar dataKey="SwapFree" fill="#bc5090" />
          <Bar dataKey="SwapTotal" fill="#ff6361" />
        </BarChart>
      </ResponsiveContainer>
    </Div>
  );
}
