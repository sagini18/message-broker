import React from "react";
import NavBar from "../components/NavBar";
import DataTable from "../components/Table";
import BasicLineChart from "../components/Chart";
import { Box, Paper } from "@mui/material";
import GraphCard from "../components/Card";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  startChannSumConnection,
  stopChannSumConnection,
} from "../store/sse/channels/sse.Thunks";
import { useMetrics } from "../store/metrics/useMetrics";

export default function Dashboard() {
  const dispatch = useDispatch();
  const { channSumEvents } = useSelector((state) => state.channSum);
  const { channelsEvents, requestsEvents, consumersEvents, messagesEvents } =
    useMetrics();

  useEffect(() => {
    dispatch(startChannSumConnection());

    return () => {
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
            dataset={messagesEvents}
          />
        </Paper>
        <Box
          display={"flex"}
          flexDirection={"column"}
          justifyContent={"space-around"}>
          <GraphCard
            title={
              "Requests for Sending Messages to Consumer Since Server Startup"
            }
            name={"No of requests"}
            color={"#5E35B1"}
            dataset={requestsEvents}
          />
          <GraphCard
            title={"Channel Creation and Removal Times Since Server Startup"}
            name={"No of channels"}
            color={"#35B175"}
            dataset={channelsEvents}
          />
          <GraphCard
            title={
              "Consumer Connection and Disconnection Times Since Server Startup"
            }
            name={"No of consumers"}
            color={"#1E88E5"}
            dataset={consumersEvents}
          />
        </Box>
      </Box>
    </div>
  );
}
