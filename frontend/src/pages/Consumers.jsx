import ChartContainer from "../components/ChartContainer";

function Consumers({consumersEvents}) {

  return (
    <ChartContainer
      title={"Consumers Cache Timeline"}
      description={
        "Tracking the consumers arrive in and are removed from the cache since the server started"
      }
      data={consumersEvents}
      bgColor={"#0099ff"}
    />
  );
}

export default Consumers;
