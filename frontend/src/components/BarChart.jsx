import * as React from 'react';
import { BarChart } from '@mui/x-charts/BarChart';
import { axisClasses } from '@mui/x-charts/ChartsAxis';

const valueFormatter = (value) => `${value}mm`;

const chartSetting = {
  yAxis: [
    {
      label: 'no of consumers',
    },
  ],
  series: [{ dataKey: 'seoul', valueFormatter }],
  height: 300,
  width: 500,
  sx: {
    [`& .${axisClasses.directionY} .${axisClasses.label}`]: {
      transform: 'translateX(-10px)',
    },
  },
};

export default function BasicBarChart({color,dataset}) {
  const tickPlacement = 'middle';
  const tickLabelPlacement = 'middle';

  return (
    <div style={{ width: '100%' }}>
      <BarChart
        dataset={dataset}
        colors={[color]}
        xAxis={[
          { scaleType: 'band', dataKey: 'month', tickPlacement, tickLabelPlacement },
        ]}
        {...chartSetting}
      />
    </div>
  );
}
