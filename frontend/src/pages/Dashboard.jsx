import React from "react";
import NavBar from "../components/NavBar";
import DataTable from "../components/Table";
import BasicLineChart from "../components/Chart";
import { Box, Paper } from "@mui/material";
import GraphCard from "../components/Card";

export default function Dashboard() {
  const rows = [
    {
      id: 1,
      channelName: "science",
      noOfMessagesInQueue: 34,
      noOfConsumers: 35,
      noOfProdcuers: 5,
      noOfMessagesInPersistence:2,
      failedMessages: 1,
    },
    {
      id: 2,
      channelName: "channel_6",
      noOfMessagesInQueue: 87,
      noOfConsumers: 42,
      noOfProdcuers: 25,
      noOfMessagesInPersistence:2,
      failedMessages: 0,
    },
    {
      id: 3,
      channelName: "channel_2",
      noOfMessagesInQueue: 34,
      noOfConsumers: 45,
      noOfProdcuers: 15,
      noOfMessagesInPersistence:2,
      failedMessages: 2,
    },
    {
      id: 4,
      channelName: "channel_7",
      noOfMessagesInQueue: 45,
      noOfConsumers: 16,
      noOfProdcuers: 5,
      noOfMessagesInPersistence:2,
      failedMessages: 0,
    },
    {
      id: 5,
      channelName: "channel_12",
      noOfMessagesInQueue: 34,
      noOfConsumers: 7,
      noOfProdcuers: 10,
      noOfMessagesInPersistence:2,
      failedMessages: 0,
    },
    {
      id: 6,
      channelName: "maths",
      noOfMessagesInQueue: 1,
      noOfConsumers: 150,
      noOfProdcuers: 5,
      noOfMessagesInPersistence:2,
      failedMessages: 1,
    },
    {
      id: 7,
      channelName: "70",
      noOfMessagesInQueue: 56,
      noOfConsumers: 44,
      noOfProdcuers: 5,
      noOfMessagesInPersistence:2,
      failedMessages: 3,
    },
    {
      id: 8,
      channelName: "data_87",
      noOfMessagesInQueue: 232,
      noOfConsumers: 36,
      noOfProdcuers: 5,
      noOfMessagesInPersistence:2,
      failedMessages: 0,
    },
    {
      id: 9,
      channelName: "channel_76",
      noOfMessagesInQueue: 12,
      noOfConsumers: 65,
      noOfProdcuers: 5,
      noOfMessagesInPersistence:2,
      failedMessages: 1,
    },
  ];

  return (
    <div>
      <NavBar />
      <Box display={"flex"} paddingInline={2} justifyContent={"space-between"}>
        <DataTable rows={rows} />
        <Box
          display={"flex"}
          flexDirection={"column"}
          justifyContent={"space-around"}>
          <GraphCard count={100} name={"No of producers"} color={"#5E35B1"} />
          <GraphCard count={800} name={"No of consumers"} color={"#35B175"} />
          <GraphCard count={14} name={"No of channels"} color={"#1E88E5"} />
        </Box>
      </Box>
      <Box display={"flex"} justifyContent={"space-around"}pt={1}>
        <Paper elevation={3} sx={{backgroundColor:"#F1F1F1"}}>
          <BasicLineChart color={"blue"} name="No of requests" />
        </Paper>
        <Paper elevation={3} sx={{backgroundColor:"#F1F1F1"}}>
          <BasicLineChart color={"blue"} name="No of messages" />
        </Paper>
      </Box>
    </div>
  );
}
