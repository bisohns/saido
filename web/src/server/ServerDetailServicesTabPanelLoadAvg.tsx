import React from "react";
import styled from "@emotion/styled";
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
  LoadingAvgData,
  ServerResponseType,
  ServerServiceNameType,
} from "./ServerType";

interface ServerDetailServicesTabPanelLoadAvgType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<LoadingAvgData>;
}

const Div = styled.div`
  margin-top: 2rem;
`;

export default function ServerDetailServicesTabPanelLoadAvg(
  props: ServerDetailServicesTabPanelLoadAvgType
) {
  const {
    serverData: {
      Message: { Data },
    },
  } = props;

  return (
    <Div>
      <ResponsiveContainer width="100%" height={500}>
        <BarChart width={900} height={500} data={[Data]}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="Name" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey="Load1M" fill="#8884d8" />
          <Bar dataKey="Load5M" fill="red" />
          <Bar dataKey="Load15M" fill="green" />
        </BarChart>
      </ResponsiveContainer>
    </Div>
  );
}
