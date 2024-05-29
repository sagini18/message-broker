import { createSlice } from "@reduxjs/toolkit";

const consumerSlice = createSlice({
  name: "consumer",
  initialState: {
    connected: false,
    consumerEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectConsumer: (state) => {
      state.connected = true;
    },
    disconnectConsumer: (state) => {
      state.connected = false;
      state.eventSourceUrl = null;
    },
    setConsumerEvents: (state, action) => {
      state.consumerEvents = action.payload;
    },
    setEventSourceUrl: (state, action) => {
      state.eventSourceUrl = action.payload;
    },
  },
});

export const {
  connectConsumer,
  disconnectConsumer,
  setConsumerEvents,
  setEventSourceUrl,
} = consumerSlice.actions;

export default consumerSlice.reducer;
