import { createAsyncThunk } from "@reduxjs/toolkit";
import {
  connectConsumer,
  disconnectConsumer,
  setConsumerEvents,
  setEventSourceUrl,
} from "./consumerSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";

export const startConsumerConnection = createAsyncThunk(
  "consumer/startConsumerConnection",
  async (_, { dispatch, getState }) => {
    const { consumer } = getState();

    if (consumer.connected) {
      return;
    }
    const CONSUMER_EVENT_SOURCE_KEY = 'consumer';
    const CONSUMER_EVENT_SOURCE_URL = "http://localhost:8080/api/v1/consumers/events";

    startEventSource(
      CONSUMER_EVENT_SOURCE_KEY,
      CONSUMER_EVENT_SOURCE_URL,
      () => {
        dispatch(connectConsumer());
        dispatch(setEventSourceUrl(CONSUMER_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        dispatch(setConsumerEvents(data));
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
