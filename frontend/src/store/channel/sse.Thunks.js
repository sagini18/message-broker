import { createAsyncThunk } from "@reduxjs/toolkit";
import {
  connectChannel,
  disconnectChannel,
  setChannelEvents,
  setEventSourceUrl,
} from "./channelSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";
import { dynamicRateLimiter } from '../../utils/dynamicRateLimiter.js';

const handleChannelEvents = (dispatch, data) => {
  dispatch(setChannelEvents(data));
};

export const startChannelConnection = createAsyncThunk(
  "channel/startChannelConnection",
  async (_, { dispatch, getState }) => {
    const { channel } = getState();

    if (channel.connected) {
      return;
    }

    const CHANNEL_EVENT_SOURCE_KEY = "channel";
    const CHANNEL_EVENT_SOURCE_URL = "http://localhost:8080/api/channel/count";


    startEventSource(
      CHANNEL_EVENT_SOURCE_KEY,
      CHANNEL_EVENT_SOURCE_URL,
      () => {
        dispatch(connectChannel());
        dispatch(setEventSourceUrl(CHANNEL_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        const load = data.length;  // Assuming data length can be used as a proxy for load
        const handleEvents = dynamicRateLimiter(handleChannelEvents, load);
        handleEvents(dispatch, data);
      },
      () => {
        dispatch(disconnectChannel());
        setTimeout(() => {
          dispatch(startChannelConnection());
        }, 1000);
      }
    );
  }
);

export const stopChannelConnection = createAsyncThunk(
  "channel/stopChannelConnection",
  async (_, { dispatch }) => {
    closeEventSource("channel");
    dispatch(disconnectChannel());
  }
);
