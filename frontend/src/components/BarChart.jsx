import * as React from "react";
import { BarChart } from "@mui/x-charts/BarChart";
import { axisClasses } from "@mui/x-charts/ChartsAxis";
import dayjs from "dayjs";

export default function BasicBarChart({ color, dataset }) {
  const tickPlacement = "middle";
  const tickLabelPlacement = "middle";

  const valueFormatter = (value) =>
    value === 0 || value === 1 ? `${value} channel` : `${value} channels`;

  let previousDate = null;

  const timestampFormatter = (timestamp) => {
    const currentTimestamp = dayjs(timestamp);
    const currentDate = currentTimestamp.format("YYYY-MM-DD");
    const formattedTimestamp =
      currentDate === previousDate
        ? currentTimestamp.format("HH:mm:ss")
        : currentTimestamp.format("YYYY-MM-DD HH:mm:ss");

    previousDate = currentDate;
    return formattedTimestamp;
  };
  
  const chartSetting = {
    yAxis: [
      {
        label: "no of channels",
        max: Math.max(...dataset.map((data) => data.Count)) + 1,
      },
    ],
    series: [{ dataKey: "Count", valueFormatter }],
    height: 300,
    width: 500,
    sx: {
      [`& .${axisClasses.directionY} .${axisClasses.label}`]: {
        transform: "translateX(-10px)",
      },
    },
  };

  console.log("dataset : ", dataset);

  return (
    <div style={{ width: "100%" }}>
      <BarChart
        dataset={dataset}
        colors={[color]}
        xAxis={[
          {
            scaleType: "band",
            dataKey: "Timestamp",
            tickPlacement,
            tickLabelPlacement,
            valueFormatter: timestampFormatter,
          },
        ]}
        {...chartSetting}
      />
    </div>
  );
}
