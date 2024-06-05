import * as React from "react";
import { BarChart } from "@mui/x-charts/BarChart";
import { axisClasses } from "@mui/x-charts/ChartsAxis";
import dayjs from "dayjs";

export default function BasicBarChart({ color, dataset }) {
  const tickPlacement = "middle";
  const tickLabelPlacement = "middle";

  const valueFormatter = (value) =>
    value === 0 || value === 1 ? `${value} channel` : `${value} channels`;

  const chartSetting = {
    yAxis: [
      {
        label: "no of channels",
      },
    ],
    series: [{ dataKey: "count", valueFormatter }],
    height: 300,
    width: 500,
    sx: {
      [`& .${axisClasses.directionY} .${axisClasses.label}`]: {
        transform: "translateX(-10px)",
      },
    },
  };

  return (
    <div style={{ width: "100%" }}>
      <BarChart
        dataset={dataset}
        colors={[color]}
        xAxis={[
          {
            scaleType: "band",
            dataKey: "time",
            tickPlacement,
            tickLabelPlacement,
          },
        ]}
        {...chartSetting}
      />
    </div>
  );
}
