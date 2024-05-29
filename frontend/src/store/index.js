import { configureStore } from '@reduxjs/toolkit';
import channelReducer from './channel/channelSlice';
import consumerReducer from './consumer/consumerSlice';
import requestReducer from "./request/requestSlice"

const store = configureStore({
  reducer: {
    channel: channelReducer,
    consumer: consumerReducer,
    request: requestReducer,
  },
});

export default store;
