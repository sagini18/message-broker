import React from "react";
import BasicLineChart from "./Chart";
import { Box, Paper, Typography } from "@mui/material";

function ChartContainer({ title, description, data, bgColor }) {
  return (
    <Box
      display={"flex"}
      flexDirection={"column"}
      paddingInline={2}
      justifyContent={"center"}
      alignItems={"center"}>
      <Box
        display={"flex"}
        flexDirection={"column"}
        alignItems={"start"}
        width={"90%"}
        >
        <Typography
          color="black"
          fontSize={"19px"}
          paddingTop={2}
          fontWeight={"700"}>
          {title}
        </Typography>
        <Typography
          color="black"
          pb={2}
          fontSize={"19px"}
          fontWeight={"500"}>
          {description}
        </Typography>
      </Box>
      <Paper
        elevation={3}
        sx={{
          backgroundColor: "#27272a",
          width: "90%",
          height: "auto",
        }}>
        <BasicLineChart color={bgColor} name="No of messages" dataset={data} />
      </Paper>
    </Box>
  );
}

export default ChartContainer;
