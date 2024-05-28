import { createAsyncThunk } from '@reduxjs/toolkit';
import { connectChannel, disconnectChannel, setEvents, setEventSource } from './channelSlice.js';

export const startSSEConnection = createAsyncThunk(
  'sse/startSSEConnection',
  async (_, { dispatch, getState }) => {
    const { channel } = getState();

    if (channel.connected) {
      return;
    }

    const eventSource = new EventSource('http://localhost:8080/api/channel/count');

    eventSource.onopen = () => {
      dispatch(connectChannel());
      dispatch(setEventSource(eventSource));
    };

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      dispatch(setEvents(data));
    };

    eventSource.onerror = () => {
      dispatch(disconnectChannel());
      setTimeout(() => {
        dispatch(startSSEConnection());
      }, 1000); 
    };
  }
);
