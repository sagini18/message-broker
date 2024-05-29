import { createAsyncThunk } from "@reduxjs/toolkit";
import {
  connectRequest,
  disconnectRequest,
  setRequestEvents,
  setEventSourceUrl,
} from "../request/requestSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";

export const startRequestConnection = createAsyncThunk(
  "request/startRequestConnection",
  async (_, { dispatch, getState }) => {
    const { request } = getState();

    if (request.connected) {
      return;
    }
    const CONSUMER_EVENT_SOURCE_KEY = 'request';
    const CONSUMER_EVENT_SOURCE_URL = "http://localhost:8080/api/request/count";

    startEventSource(
      CONSUMER_EVENT_SOURCE_KEY,
      CONSUMER_EVENT_SOURCE_URL,
      () => {
        dispatch(connectRequest());
        dispatch(setEventSourceUrl(CONSUMER_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        dispatch(setRequestEvents(data));
      },
      () => {
        dispatch(disconnectRequest());
        setTimeout(() => {
          dispatch(startRequestConnection());
        }, 1000);
      }
    );
  }
);

export const stopRequestConnection = createAsyncThunk(
  "request/stopRequestConnection",
  async (_, { dispatch }) => {
    closeEventSource("request");
    dispatch(disconnectRequest());
  }
);
