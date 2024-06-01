import { createAsyncThunk } from "@reduxjs/toolkit";
import {
  connectRequest,
  disconnectRequest,
  setRequestEvents,
  setEventSourceUrl,
} from "../request/requestSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";
import { dynamicRateLimiter } from '../../utils/dynamicRateLimiter.js';

const handleRequestEvents = (dispatch, data) => {
  dispatch(setRequestEvents(data));
};


export const startRequestConnection = createAsyncThunk(
  "request/startRequestConnection",
  async (_, { dispatch, getState }) => {
    const { request } = getState();

    if (request.connected) {
      return;
    }
    const REQ_EVENT_SOURCE_KEY = 'request';
    const REQ_EVENT_SOURCE_URL = "http://localhost:8080/api/request/count";

    startEventSource(
      REQ_EVENT_SOURCE_KEY,
      REQ_EVENT_SOURCE_URL,
      () => {
        dispatch(connectRequest());
        dispatch(setEventSourceUrl(REQ_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        const load = data.length;  // Assuming data length can be used as a proxy for load
        const handleEvents = dynamicRateLimiter(handleRequestEvents, load);
        handleEvents(dispatch, data);
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
