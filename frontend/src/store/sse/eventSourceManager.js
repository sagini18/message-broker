let eventSources = {};

export const startEventSource = (key, url, onOpen, onMessage, onError) => {
  if (!eventSources[key]) {
    eventSources[key] = new EventSource(url);
  }

  eventSources[key].onopen = onOpen;
  eventSources[key].onmessage = onMessage;
  eventSources[key].onerror = onError;
};

export const closeEventSource = (key) => {
  if (eventSources[key]) {
    eventSources[key].close();
    delete eventSources[key];
  }
};

export const getEventSource = (key) => eventSources[key];
