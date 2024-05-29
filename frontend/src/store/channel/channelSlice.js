import { createSlice } from '@reduxjs/toolkit';

const channelSlice = createSlice({
  name: 'channel',
  initialState: {
    connected: false,
    channelEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectChannel: (state) => {
      state.connected = true;
    },
    disconnectChannel: (state) => {
      state.connected = false;
      state.eventSourceUrl = null;
    },
    setChannelEvents: (state, action) => {
      state.channelEvents = action.payload;
    },
    setEventSourceUrl: (state, action) => {
      state.eventSourceUrl = action.payload;
    },
  },
});

export const { connectChannel, disconnectChannel, setChannelEvents, setEventSourceUrl } = channelSlice.actions;
export default channelSlice.reducer;
