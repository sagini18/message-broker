import React from "react";
import NavBar from "../components/NavBar";
import DataTable from "../components/Table";
import BasicLineChart from "../components/Chart";
import { Box, Paper } from "@mui/material";
import GraphCard from "../components/Card";
import { useEffect } from "react";
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
import {
  startChannSumConnection,
  stopChannSumConnection,
} from "../store/channel_summary/sse.Thunks";

export default function Dashboard() {
  const dispatch = useDispatch();
  const { channelEvents } = useSelector((state) => state.channel);
  const { consumerEvents } = useSelector((state) => state.consumer);
  const { requestEvents } = useSelector((state) => state.request);
  const { msgEvents } = useSelector((state) => state.message);
  const { channSumEvents } = useSelector((state) => state.channSum);

  useEffect(() => {
    dispatch(startChannelConnection());
    dispatch(startConsumerConnection());
    dispatch(startRequestConnection());
    dispatch(startMsgConnection());
    dispatch(startChannSumConnection());

    return () => {
      dispatch(stopChannelConnection());
      dispatch(stopConsumerConnection());
      dispatch(stopRequestConnection());
      dispatch(stopMsgConnection());
      dispatch(stopChannSumConnection());
    };
  }, [dispatch]);

  return (
    <div>
      <NavBar />
      <Box display={"flex"} paddingInline={2} justifyContent={"space-evenly"}>
        <DataTable rows={channSumEvents} />
      </Box>
      <Box display={"flex"} justifyContent={"space-evenly"} pt={1}>
        <Paper elevation={3} sx={{ backgroundColor: "#F1F1F1", width:"45vw", display:"flex",justifyContent:"center" }}>
          <BasicLineChart
            color={"blue"}
            name="No of messages"
            dataset={msgEvents}
          />
        </Paper>
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
    </div>
  );
}
