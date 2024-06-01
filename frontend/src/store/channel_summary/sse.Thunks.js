import { createAsyncThunk } from "@reduxjs/toolkit";
import { connectChannSum, disconnectChannSum, setChannSumEvents, setEventSourceUrl } from "./channSumSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";
import { dynamicRateLimiter } from '../../utils/dynamicRateLimiter.js';

const handleChannSumEvents = (dispatch, data) => {
  dispatch(setChannSumEvents(data));
};


export const startChannSumConnection = createAsyncThunk(
  "channSum/startChannSumConnection",
  async (_, { dispatch, getState }) => {
    const { channSum } = getState();

    if (channSum.connected) {
      return;
    }

    const CHANN_SUM_EVENT_SOURCE_KEY = "channSum";
    const CHANN_SUM_EVENT_SOURCE_URL = "http://localhost:8080/api/channel/all";


    startEventSource(
      CHANN_SUM_EVENT_SOURCE_KEY,
      CHANN_SUM_EVENT_SOURCE_URL,
      () => {
        dispatch(connectChannSum());
        dispatch(setEventSourceUrl(CHANN_SUM_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        const load = data.length;
        const handleEvents = dynamicRateLimiter(handleChannSumEvents, load);
        handleEvents(dispatch, data);
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
