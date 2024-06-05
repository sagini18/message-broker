import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import NavBar from "./components/NavBar";
import ChannelsTable from "./pages/ChannelsTable";
import Messages from "./pages/Messages";
import Requests from "./pages/Requestes";
import Consumers from "./pages/Consumers";
import Channels from "./pages/Channels";
import { useMetrics } from "./store/metrics/useMetrics";

function App() {
  const {channelsEvents, requestsEvents,messagesEvents,consumersEvents} = useMetrics();
  return (
    <Router>
    <NavBar />
    <Routes>
      <Route path="/table" element={<ChannelsTable />} />
      <Route path="/message" element={<Messages messagesEvents={messagesEvents}/>} />
      <Route path="/request" element={<Requests requestsEvents={requestsEvents} />} />
      <Route path="/consumer" element={<Consumers consumersEvents={consumersEvents} />} />
      <Route path="/channel" element={<Channels channelsEvents={channelsEvents} />} />
      <Route path="/" element={<ChannelsTable />} /> {/* Default route */}
    </Routes>
  </Router>
  );
}

export default App;
