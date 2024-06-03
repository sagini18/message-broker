import { createSlice } from '@reduxjs/toolkit';

const messageSlice = createSlice({
  name: 'message',
  initialState: {
    msgConnected: false,
    msgEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectMsg: (state) => {
      state.msgConnected = true;
    },
    disconnectMsg: (state) => {
      state.msgConnected = false;
      state.eventSourceUrl = null;
    },
    setMsgEvents: (state, action) => {
      state.msgEvents = action.payload;
    },
    setEventSourceUrl: (state, action) => {
      state.eventSourceUrl = action.payload;
    },
  },
});

export const { connectMsg, disconnectMsg, setMsgEvents, setEventSourceUrl } = messageSlice.actions;
export default messageSlice.reducer;
