import { createAsyncThunk } from "@reduxjs/toolkit";
import { connectChannSum, disconnectChannSum, setChannSumEvents, setEventSourceUrl } from "./channSumSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";

export const startChannSumConnection = createAsyncThunk(
  "channSum/startChannSumConnection",
  async (_, { dispatch, getState }) => {
    const { channSum } = getState();

    if (channSum.connected) {
      return;
    }

    const CHANN_SUM_EVENT_SOURCE_KEY = "channSum";
    const CHANN_SUM_EVENT_SOURCE_URL = "http://localhost:8080/api/v1/channels";

    startEventSource(
      CHANN_SUM_EVENT_SOURCE_KEY,
      CHANN_SUM_EVENT_SOURCE_URL,
      () => {
        dispatch(connectChannSum());
        dispatch(setEventSourceUrl(CHANN_SUM_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        dispatch(setChannSumEvents(data));
      },
      () => {
        dispatch(disconnectChannSum());
        setTimeout(() => {
          dispatch(startChannSumConnection());
        }, 1000);
      }
    );
  }
);

export const stopChannSumConnection = createAsyncThunk(
  "channSum/stopChannSumConnection",
  async (_, { dispatch }) => {
    closeEventSource("channSum");
    dispatch(disconnectChannSum());
  }
);
