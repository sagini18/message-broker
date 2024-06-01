import throttle from 'lodash/throttle';
import debounce from 'lodash/debounce';

const THROTTLE_INTERVAL = 1000;  // Default throttle interval
const DEBOUNCE_INTERVAL = 100;   // Default debounce interval

let currentThrottleInterval = THROTTLE_INTERVAL;
let currentDebounceInterval = DEBOUNCE_INTERVAL;

const adjustRate = (load) => {
    console.log('load:', load);
  // Adjust throttle and debounce intervals based on the load
  if (load > 48) {
    currentThrottleInterval = 2000;  // Increase throttle interval under heavy load
    currentDebounceInterval = 200;   // Increase debounce interval under heavy load
  } else {
    currentThrottleInterval = THROTTLE_INTERVAL;  // Default throttle interval
    currentDebounceInterval = DEBOUNCE_INTERVAL;  // Default debounce interval
  }
};

export const dynamicRateLimiter = (fn, load) => {
  adjustRate(load);
  return debounce(throttle(fn, currentThrottleInterval), currentDebounceInterval);
};
