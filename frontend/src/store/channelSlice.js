import { createSlice } from '@reduxjs/toolkit';

const channelSlice = createSlice({
  name: 'channel',
  initialState: {
    connected: false,
    events: [],
    eventSource: null,
  },
  reducers: {
    connectChannel: (state) => {
      state.connected = true;
    },
    disconnectChannel: (state) => {
      state.connected = false;
      if (state.eventSource) {
        state.eventSource.close();
        state.eventSource = null;
      }
    },
    setEvents: (state, action) => {
      state.events = action.payload;
    },
    setEventSource: (state, action) => {
      state.eventSource = action.payload;
    },
  },
});

export const { connectChannel, disconnectChannel, setEvents, setEventSource } = channelSlice.actions;
export default channelSlice.reducer;
