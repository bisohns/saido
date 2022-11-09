import React from 'react';
import styled from '@emotion/styled';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from 'recharts';

import {
  DiskData,
  ServerResponseType,
  ServerServiceNameType,
} from './ServerType';

interface ServerDetailServicesTabPanelDiskType {
  serverName: ServerServiceNameType;
  serverData: ServerResponseType<Array<DiskData>>;
}

const Div = styled.div`
  margin-top: 2rem;
`;

export default function ServerDetailServicesTabPanelDisk(
  props: ServerDetailServicesTabPanelDiskType
) {
  const {
    serverData: {
      Message: { Data },
    },
  } = props;

  return (
    <Div>
      {
        /* <ServicesTabPanel /> */
        <BarChart width={900} height={500} data={Data}>
          <CartesianGrid strokeDasharray='3 3' />
          <XAxis dataKey='FileSystem' />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey='Available' fill='#8884d8' />
          <Bar dataKey='Used' fill='red' />
          <Bar dataKey='Size' fill='green' />
        </BarChart>
      }
    </Div>
  );
}
