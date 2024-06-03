import { createSlice } from "@reduxjs/toolkit";

const consumerSlice = createSlice({
  name: "consumer",
  initialState: {
    consumerConnected: false,
    consumerEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectConsumer: (state) => {
      state.consumerConnected = true;
    },
    disconnectConsumer: (state) => {
      state.consumerConnected = false;
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
