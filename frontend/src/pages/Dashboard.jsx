import React from "react";
import NavBar from "../components/NavBar";
import DataTable from "../components/Table";
import BasicLineChart from "../components/Chart";
import { Box } from "@mui/material";

export default function Dashboard() {
  return (
    <div>
      <NavBar/>
      <Box display={"flex"} >
        <BasicLineChart />
        <BasicLineChart />
      </Box>
      <Box display={"flex"} paddingInline={2}>
      <DataTable />
      </Box>
    </div>
  );
}
