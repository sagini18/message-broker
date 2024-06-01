import { configureStore } from '@reduxjs/toolkit';
import channelReducer from './channel/channelSlice';
import consumerReducer from './consumer/consumerSlice';
import requestReducer from "./request/requestSlice"
import messageReducer from "./message/messageSlice"
import channelSummaryReducer from "./channel_summary/channSumSlice"

const store = configureStore({
  reducer: {
    channel: channelReducer,
    consumer: consumerReducer,
    request: requestReducer,
    message: messageReducer,
    channSum: channelSummaryReducer,
  },
});

export default store;
