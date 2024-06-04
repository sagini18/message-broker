import { configureStore } from '@reduxjs/toolkit';
import channelSummaryReducer from "./channels/channSumSlice"

const store = configureStore({
  reducer: {
    channSum: channelSummaryReducer,
  },
});

export default store;
