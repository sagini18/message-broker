import { createSlice } from "@reduxjs/toolkit";

const requestSlice = createSlice({
  name: "request",
  initialState: {
    reqConnected: false,
    requestEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectRequest: (state) => {
      state.reqConnected = true;
    },
    disconnectRequest: (state) => {
      state.reqConnected = false;
      state.eventSourceUrl = null;
    },
    setRequestEvents: (state, action) => {
      state.requestEvents = action.payload;
    },
    setEventSourceUrl: (state, action) => {
      state.eventSourceUrl = action.payload;
    },
  },
});

export const {
  connectRequest,
  disconnectRequest,
  setRequestEvents,
  setEventSourceUrl,
} = requestSlice.actions;

export default requestSlice.reducer;
