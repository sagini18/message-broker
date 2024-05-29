import * as React from "react";
import { LineChart } from "@mui/x-charts/LineChart";
import dayjs from "dayjs";

export default function BasicLineChart({ color, name, dataset }) {
  const xAxisData = dataset?.map((item) => dayjs(item?.Timestamp)?.toDate());
  const seriesData = dataset?.map((item) => item?.Count);

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

  return (
    <LineChart
      colors={[color]}
      xAxis={[
        {
          scaleType: "point",
          data: xAxisData,
          valueFormatter: timestampFormatter,
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
      width={500}
      height={300}
      slotProps={{
        legend: {
          hidden : !(dataset?.length > 0)
        }
      }}
    />
  );
}
