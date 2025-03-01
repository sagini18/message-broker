import axios from "axios";
import { useEffect, useState } from "react";

export const useMetrics = () => {
  const [channelsEvents, setChannelsEvents] = useState([]);
  const [requestsEvents, setRequestsEvents] = useState([]);
  const [consumersEvents, setConsumersEvents] = useState([]);
  const [messagesEvents, setMessagesEvents] = useState([]);
  
  useEffect(() => {
    const fetchData = async () => {
      const metricsToFetch = [
        'messages_events',
        'consumers_events',
        'requests_events',
        'channels_events',
    ];
      try {
        let response
        metricsToFetch.map(async(metric)=>{
          response = await axios.get(
            `http://localhost:9090/api/v1/query?query=${metric}`
          );
          
          const result = response.data?.data?.result?.[0];
          if (result.metric.__name__ === metric) {
            const newData = {
              time: new Date().toLocaleTimeString(),
              count: result.value[1],
            };
  
            switch (metric) {
              case 'channels_events':
                setChannelsEvents((prevData) => [...prevData, newData]);
                break;
              case 'requests_events':
                setRequestsEvents((prevData) => [...prevData, newData]);
                break;
              case 'consumers_events':
                setConsumersEvents((prevData) => [...prevData, newData]);
                break;
              case 'messages_events':
                setMessagesEvents((prevData) => [...prevData, newData]);
                break;
              default:
                break;
            }
        }
      })
      } catch (error) {
        console.error("Error fetching metrics:", error);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, [channelsEvents, requestsEvents, consumersEvents, messagesEvents]);

  return { channelsEvents, requestsEvents, consumersEvents, messagesEvents };
};
