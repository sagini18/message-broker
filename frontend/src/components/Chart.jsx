import React, { useEffect, useState } from "react";
import { LineChart } from "@mui/x-charts/LineChart";

export default function BasicLineChart({ color, name, dataset }) {
  const [chartWidth, setChartWidth] = useState(window.innerWidth * 0.88);
  const [chartHeight, setChartHeight] = useState(500);

  useEffect(() => {
    const handleResize = () => {
      setChartWidth(window.innerWidth * 0.88);
      setChartHeight(350);
    };
    window.addEventListener("resize", handleResize);

    const style = document.createElement("style");
    style.innerHTML = `
    .line-chart text {
      fill: #adbdb6; 
      font-family: Arial, sans-serif; 
    }
  `;
    document.head.appendChild(style);
    return () => {
      window.removeEventListener("resize", handleResize);
      document.head.removeChild(style);
    };
  }, []);

  const xAxisData = dataset?.map((item) => item?.time);
  const seriesData = dataset?.map((item) => item?.count);

  return (
    <div className="line-chart">
      <LineChart
        colors={[color]}
        tooltip={{
          slots: {},
        }}
        margin={{ top: 50, right: 50, left: 50, bottom: 50 }}
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
        className="line-chart"
        width={chartWidth}
        height={chartHeight}
        slotProps={{
          legend: {
            hidden: !(dataset?.length > 0),
            labelStyle: {
              fill: '#adbdb6',
            }
          },
        }}
        sx={{
          "& .MuiChartsAxis-tickContainer .MuiChartsAxis-tickLabel": {
            strokeWidth: "0.7",
            fill: "#adbdb6",
            fontFamily: "Arial, sans-serif",
          },
          "& .MuiChartsAxis-bottom .MuiChartsAxis-line": {
            stroke: "#adbdb6",
            strokeWidth: 0.5,
          },
          "& .MuiChartsAxis-left .MuiChartsAxis-line": {
            stroke: "#adbdb6",
            strokeWidth: 0.5,
          },
          "& .MuiChartsAxis-tickContainer .MuiChartsAxis-tick": {
            stroke: "#adbdb6",
            strokeWidth: "0.5",
          },
          "& .MuiChartsAxisHighlight-root":{
            stroke: "#adbdb6",
            strokeWidth:0.5,
          },
        }}
      />
    </div>
  );
}
