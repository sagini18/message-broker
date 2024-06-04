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
} from "../store/channels_events/sse.Thunks";
import {
  startConsumerConnection,
  stopConsumerConnection,
} from "../store/consumers_events/sse.Thunks";
import {
  startRequestConnection,
  stopRequestConnection,
} from "../store/requests_events/sse.Thunks";
import {
  startMsgConnection,
  stopMsgConnection,
} from "../store/messages_events/sse.Thunks";
import {
  startChannSumConnection,
  stopChannSumConnection,
} from "../store/channels/sse.Thunks";

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
        <Paper
          elevation={3}
          sx={{
            backgroundColor: "#F1F1F1",
            width: "45vw",
            display: "flex",
            justifyContent: "center",
          }}>
          {/* "Message Arrival and Removal Times from Cache Since Server Startup" */}
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
            title={"Requests for Sending Messages to Consumer Since Server Startup"}
            name={"No of requests"}
            color={"#5E35B1"}
            dataset={requestEvents}
          />
          <GraphCard
            title={"Channel Creation and Removal Times Since Server Startup"}
            name={"No of channels"}
            color={"#35B175"}
            dataset={channelEvents}
          />
          <GraphCard
            title={"Consumer Connection and Disconnection Times Since Server Startup"}
            name={"No of consumers"}
            color={"#1E88E5"}
            dataset={consumerEvents}
          />
        </Box>
      </Box>
    </div>
  );
}
