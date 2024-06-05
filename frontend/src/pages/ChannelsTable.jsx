import { Box } from "@mui/material";
import { useDispatch, useSelector } from "react-redux";
import { useEffect } from "react";
import {
  startChannSumConnection,
  stopChannSumConnection,
} from "../store/sse/channels/sse.Thunks";
import DataTable from "../components/Table";
import { Typography } from "@mui/material";

function ChannelsTable() {
  const dispatch = useDispatch();
  const { channSumEvents } = useSelector((state) => state.channSum);

  useEffect(() => {
    dispatch(startChannSumConnection());

    return () => {
      dispatch(stopChannSumConnection());
    };
  }, [dispatch]);

  return (
    <Box
      display={"flex"}
      paddingInline={2}
      justifyContent={"space-evenly"}
      height={"90vh"}>
      <Box>
        <Typography
          color="black"
          fontSize={"19px"}
          paddingTop={2}
          paddingBottom={1}
          fontWeight={"700"}>
          Channels In-Cache Summary
        </Typography>
        <DataTable rows={channSumEvents} />
      </Box>
    </Box>
  );
}

export default ChannelsTable;
