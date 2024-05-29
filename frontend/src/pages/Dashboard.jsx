import React from "react";
import NavBar from "../components/NavBar";
import DataTable from "../components/Table";
import BasicLineChart from "../components/Chart";
import { Box, Paper } from "@mui/material";
import GraphCard from "../components/Card";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  startChannelConnection,
  stopChannelConnection,
} from "../store/channel/sse.Thunks";
import {
  startConsumerConnection,
  stopConsumerConnection,
} from "../store/consumer/sse.Thunks";
import {
  startRequestConnection,
  stopRequestConnection,
} from "../store/request/sse.Thunks";
import {
  startMsgConnection,
  stopMsgConnection,
} from "../store/message/sse.Thunks";

export default function Dashboard() {
  const dispatch = useDispatch();
  const { channelEvents } = useSelector((state) => state.channel);
  const { consumerEvents } = useSelector((state) => state.consumer);
  const { requestEvents } = useSelector((state) => state.request);
  const { msgEvents } = useSelector((state) => state.message);

  useEffect(() => {
    dispatch(startChannelConnection());
    dispatch(startConsumerConnection());
    dispatch(startRequestConnection());
    dispatch(startMsgConnection());

    return () => {
      dispatch(stopChannelConnection());
      dispatch(stopConsumerConnection());
      dispatch(stopRequestConnection());
      dispatch(stopMsgConnection());
    };
  }, [dispatch]);

  const rows = [
    {
      id: 1,
      channelName: "science",
      noOfMessagesInQueue: 34,
      noOfConsumers: 35,
      noOfRequests: 5,
      noOfMessagesInPersistence: 2,
      failedMessages: 1,
    },
    {
      id: 2,
      channelName: "channel_6",
      noOfMessagesInQueue: 87,
      noOfConsumers: 42,
      noOfRequests: 25,
      noOfMessagesInPersistence: 2,
      failedMessages: 0,
    },
    {
      id: 3,
      channelName: "channel_2",
      noOfMessagesInQueue: 34,
      noOfConsumers: 45,
      noOfRequests: 15,
      noOfMessagesInPersistence: 2,
      failedMessages: 2,
    },
    {
      id: 4,
      channelName: "channel_7",
      noOfMessagesInQueue: 45,
      noOfConsumers: 16,
      noOfRequests: 5,
      noOfMessagesInPersistence: 2,
      failedMessages: 0,
    },
    {
      id: 5,
      channelName: "channel_12",
      noOfMessagesInQueue: 34,
      noOfConsumers: 7,
      noOfRequests: 10,
      noOfMessagesInPersistence: 2,
      failedMessages: 0,
    },
    {
      id: 6,
      channelName: "maths",
      noOfMessagesInQueue: 1,
      noOfConsumers: 150,
      noOfRequests: 5,
      noOfMessagesInPersistence: 2,
      failedMessages: 1,
    },
    {
      id: 7,
      channelName: "70",
      noOfMessagesInQueue: 56,
      noOfConsumers: 44,
      noOfRequests: 5,
      noOfMessagesInPersistence: 2,
      failedMessages: 3,
    },
    {
      id: 8,
      channelName: "data_87",
      noOfMessagesInQueue: 232,
      noOfConsumers: 36,
      noOfRequests: 5,
      noOfMessagesInPersistence: 2,
      failedMessages: 0,
    },
    {
      id: 9,
      channelName: "channel_76",
      noOfMessagesInQueue: 12,
      noOfConsumers: 65,
      noOfRequests: 5,
      noOfMessagesInPersistence: 2,
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
          <GraphCard
            name={"No of requests"}
            color={"#5E35B1"}
            dataset={requestEvents}
          />
          <GraphCard
            name={"No of channels"}
            color={"#35B175"}
            dataset={channelEvents}
          />
          <GraphCard
            name={"No of consumers"}
            color={"#1E88E5"}
            dataset={consumerEvents}
          />
        </Box>
      </Box>
      <Box display={"flex"} justifyContent={"space-around"} pt={1}>
        {/* <Paper elevation={3} sx={{backgroundColor:"#F1F1F1"}}>
          <BasicLineChart color={"blue"} name="No of requests" xAxisData={xAxisData} seriesData={seriesData} />
        </Paper> */}
        <Paper elevation={3} sx={{ backgroundColor: "#F1F1F1" }}>
          <BasicLineChart
            color={"blue"}
            name="No of messages"
            dataset={msgEvents}
          />
        </Paper>
      </Box>
    </div>
  );
}
