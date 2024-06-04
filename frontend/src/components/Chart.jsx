import * as React from "react";
import { LineChart } from "@mui/x-charts/LineChart";

export default function BasicLineChart({ color, name, dataset }) {
  const xAxisData = dataset?.map((item) => item?.time);
  const seriesData = dataset?.map((item) => item?.count);

  return (
    <LineChart
      colors={[color]}
      margin={{ top: 50, right: 50, left: 75, bottom: 40 }}
      xAxis={[
        {
          scaleType: "point",
          data: xAxisData,
        },
      ]}
      yAxis={[
        {
          data: seriesData,
        },
      ]}
      series={[
        {
          data: seriesData,
          label: dataset?.length > 0 ? name : "",
        },
      ]}
      width={560}
      height={300}
      slotProps={{
        legend: {
          hidden : !(dataset?.length > 0)
        }
      }}
    />
  );
}
