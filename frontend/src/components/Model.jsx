import React from "react";
import { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { styled } from "@mui/material/styles";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import IconButton from "@mui/material/IconButton";
import CloseIcon from "@mui/icons-material/Close";
import BasicBarChart from "./BarChart";
import BasicLineChart from "./Chart";
import { startSSEConnection } from "../store/sse.Thunks";

const BootstrapDialog = styled(Dialog)(({ theme }) => ({
  "& .MuiDialogContent-root": {
    padding: theme.spacing(2),
  },
  "& .MuiDialogActions-root": {
    padding: theme.spacing(1),
  },
}));

export default function CardModel({ open, handleClose, name, color }) {
  const dispatch = useDispatch();
  const { connected, events } = useSelector((state) => state.channel);

  useEffect(() => {
    dispatch(startSSEConnection());
  }, [dispatch]);

  const xAxisData = [1, 2, 3, 5, 8, 10, 12];
  const seriesData = [2, 5.5, 2, 8.5, 1.5, 5, 9];

  return (
    <React.Fragment>
      <BootstrapDialog
        onClose={handleClose}
        aria-labelledby="customized-dialog-title"
        open={open}>
        <DialogTitle sx={{ m: 0, p: 2 }} id="customized-dialog-title">
          {name}
        </DialogTitle>
        <IconButton
          aria-label="close"
          onClick={handleClose}
          sx={{
            position: "absolute",
            right: 8,
            top: 8,
            color: (theme) => theme.palette.grey[500],
          }}>
          <CloseIcon />
        </IconButton>
        <DialogContent dividers>
          {name === "No of channels" ? (
            <BasicBarChart color={color} dataset={events} />
          ) : (
            <BasicLineChart
              color={color}
              name={name}
              xAxisData={xAxisData}
              seriesData={seriesData}
            />
          )}
        </DialogContent>
      </BootstrapDialog>
    </React.Fragment>
  );
}
