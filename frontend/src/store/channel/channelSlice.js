import { createSlice } from '@reduxjs/toolkit';

const channelSlice = createSlice({
  name: 'channel',
  initialState: {
    channConnected: false,
    channelEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectChannel: (state) => {
      state.channConnected = true;
    },
    disconnectChannel: (state) => {
      state.channConnected = false;
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
