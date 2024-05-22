import * as React from "react";
import { LineChart } from "@mui/x-charts/LineChart";

export default function BasicLineChart({ color, name, xAxisData, seriesData }) {
  console.log(xAxisData, seriesData);
  return (
    <LineChart
      colors={[color]}
      xAxis={[
        {
          // scaleType: "time",
          data: xAxisData,
        },
      ]}
      series={[
        {
          data: seriesData,
          label: `${name}`,
        },
      ]}
      width={500}
      height={300}
    />
  );
}
