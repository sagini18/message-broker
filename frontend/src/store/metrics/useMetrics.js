import axios from "axios";
import { useEffect, useState } from "react";

export const useMetrics = () => {
  const [channelsEvents, setChannelsEvents] = useState(
    JSON.parse(localStorage.getItem("channelsEvents")) || []
  );
  const [requestsEvents, setRequestsEvents] = useState(
    JSON.parse(localStorage.getItem("requestsEvents")) || []
  );
  const [consumersEvents, setConsumersEvents] = useState(
    JSON.parse(localStorage.getItem("consumersEvents")) || []
  );
  const [messagesEvents, setMessagesEvents] = useState(
    JSON.parse(localStorage.getItem("messagesEvents")) || []
  );

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(
          "http://localhost:8080/api/v1/metrics"
        );
        if (response.data) {
          const channelsMatch = response.data.match(/channels_events\s+(\d+)/);
          const requestsMatch = response.data.match(/requests_events\s+(\d+)/);
          const consumersMatch = response.data.match(/consumers_events\s+(\d+)/);
          const messagesMatch = response.data.match(/messages_events\s+(\d+)/);
          const channels = channelsMatch ? Number(channelsMatch[1]) : null;
          const requests = requestsMatch ? Number(requestsMatch[1]) : null;
          const consumers = consumersMatch ? Number(consumersMatch[1]) : null;
          const messages = messagesMatch ? Number(messagesMatch[1]) : null;

          if (channels !== null) {
            setChannelsEvents((prevData) => [
              ...prevData,
              {
                time: new Date().toLocaleTimeString(),
                count: channels,
              },
            ]);
          }

          if (requests !== null) {
            setRequestsEvents((prevData) => [
              ...prevData,
              {
                time: new Date().toLocaleTimeString(),
                count: requests,
              },
            ]);
          }

          if (consumers !== null) {
            setConsumersEvents((prevData) => [
              ...prevData,
              {
                time: new Date().toLocaleTimeString(),
                count: consumers,
              },
            ]);
          }

          if (messages !== null) {
            setMessagesEvents((prevData) => [
              ...prevData,
              {
                time: new Date().toLocaleTimeString(),
                count: messages,
              },
            ]);
          }
        }
      } catch (error) {
        console.error("Error fetching metrics:", error);
      }
    };

    const intervalId = setInterval(fetchData, 5000);

    localStorage.setItem("channelsEvents", JSON.stringify(channelsEvents));
    localStorage.setItem("requestsEvents", JSON.stringify(requestsEvents));
    localStorage.setItem("consumersEvents", JSON.stringify(consumersEvents));
    localStorage.setItem("messagesEvents", JSON.stringify(messagesEvents));

    return () => clearInterval(intervalId);
  }, [channelsEvents, requestsEvents, consumersEvents, messagesEvents]);

  return { channelsEvents, requestsEvents, consumersEvents, messagesEvents };
};
