import { configureStore } from '@reduxjs/toolkit';
import channelReducer from './channelSlice';

const store = configureStore({
  reducer: {
    channel: channelReducer,
  },
});

export default store;
