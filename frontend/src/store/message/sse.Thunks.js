import { createAsyncThunk } from "@reduxjs/toolkit";
import {
  connectMsg,
  disconnectMsg,
  setMsgEvents,
  setEventSourceUrl,
} from "../message/messageSlice.js";
import { startEventSource, closeEventSource } from "../eventSourceManager.js";

export const startMsgConnection = createAsyncThunk(
  "message/startMsgConnection",
  async (_, { dispatch, getState }) => {
    const { message } = getState();

    if (message.connected) {
      return;
    }
    const MSG_EVENT_SOURCE_KEY = 'message';
    const MSG_EVENT_SOURCE_URL = "http://localhost:8080/api/message/count";

    startEventSource(
      MSG_EVENT_SOURCE_KEY,
      MSG_EVENT_SOURCE_URL,
      () => {
        dispatch(connectMsg());
        dispatch(setEventSourceUrl(MSG_EVENT_SOURCE_URL));
      },
      (event) => {
        const data = JSON.parse(event.data);
        dispatch(setMsgEvents(data));
      },
      () => {
        dispatch(disconnectMsg());
        setTimeout(() => {
          dispatch(startMsgConnection());
        }, 1000);
      }
    );
  }
);

export const stopMsgConnection = createAsyncThunk(
  "message/stopMsgConnection",
  async (_, { dispatch }) => {
    closeEventSource("message");
    dispatch(disconnectMsg());
  }
);
