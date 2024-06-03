import { createSlice } from '@reduxjs/toolkit';

const channSumSlice = createSlice({
  name: 'channSum',
  initialState: {
    channSumConnected: false,
    channSumEvents: [],
    eventSourceUrl: null,
  },
  reducers: {
    connectChannSum: (state) => {
      state.channSumConnected = true;
    },
    disconnectChannSum: (state) => {
      state.channSumConnected = false;
      state.eventSourceUrl = null;
    },
    setChannSumEvents: (state, action) => {
      state.channSumEvents = action.payload;
    },
    setEventSourceUrl: (state, action) => {
      state.eventSourceUrl = action.payload;
    },
  },
});

export const { connectChannSum, disconnectChannSum, setChannSumEvents, setEventSourceUrl } = channSumSlice.actions;
export default channSumSlice.reducer;
