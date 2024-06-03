import { configureStore } from '@reduxjs/toolkit';
import channelReducer from './channels_events/channelSlice';
import consumerReducer from './consumers_events/consumerSlice';
import requestReducer from "./requests_events/requestSlice"
import messageReducer from "./messages_events/messageSlice"
import channelSummaryReducer from "./channels/channSumSlice"

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
