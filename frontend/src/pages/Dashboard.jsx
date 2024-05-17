import React from "react";
import NavBar from "../components/NavBar";
import DataTable from "../components/Table";
import BasicLineChart from "../components/Chart";
import { Box, Paper } from "@mui/material";
import GraphCard from "../components/Card";

export default function Dashboard() {
  return (
    <div>
      <NavBar />
      <Box display={"flex"} paddingInline={2} justifyContent={"space-between"}>
        <DataTable />
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
