import { createAsyncThunk } from "@reduxjs/toolkit";
import {
  connectConsumer,
  disconnectConsumer,
  setConsumerEvents,
  setEventSourceUrl,
} from "../consumer/consumerSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";
import { dynamicRateLimiter } from '../../utils/dynamicRateLimiter.js';

const handleConsumerEvents = (dispatch, data) => {
  dispatch(setConsumerEvents(data));
};

export const startConsumerConnection = createAsyncThunk(
  "consumer/startConsumerConnection",
  async (_, { dispatch, getState }) => {
    const { consumer } = getState();

    if (consumer.connected) {
      return;
    }
    const CONSUMER_EVENT_SOURCE_KEY = 'consumer';
    const CONSUMER_EVENT_SOURCE_URL = "http://localhost:8080/api/consumer/count";

    startEventSource(
      CONSUMER_EVENT_SOURCE_KEY,
      CONSUMER_EVENT_SOURCE_URL,
      () => {
        dispatch(connectConsumer());
        dispatch(setEventSourceUrl(CONSUMER_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        const load = data.length;
        const handleEvents = dynamicRateLimiter(handleConsumerEvents, load);
        handleEvents(dispatch, data);
      },
      () => {
        dispatch(disconnectConsumer());
        setTimeout(() => {
          dispatch(startConsumerConnection());
        }, 1000);
      }
    );
  }
);

export const stopConsumerConnection = createAsyncThunk(
  "consumer/stopConsumerConnection",
  async (_, { dispatch }) => {
    closeEventSource("consumer");
    dispatch(disconnectConsumer());
  }
);
