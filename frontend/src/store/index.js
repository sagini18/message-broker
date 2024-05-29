import { configureStore } from '@reduxjs/toolkit';
import channelReducer from './channel/channelSlice';
import consumerReducer from './consumer/consumerSlice';

const store = configureStore({
  reducer: {
    channel: channelReducer,
    consumer: consumerReducer,
  },
});

export default store;
