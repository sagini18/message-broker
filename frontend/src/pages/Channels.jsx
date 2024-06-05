import ChartContainer from "../components/ChartContainer";

function Channels({channelsEvents}) {
  return (
    <ChartContainer
      title={"Channels Cache Timeline"}
      description={
        "Tracking the channels arrive in and are removed from the cache since the server started"
      }
      data={channelsEvents}
      bgColor={"#01cb5b"}
    />
  );
}

export default Channels;
