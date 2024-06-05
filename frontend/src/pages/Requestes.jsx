import ChartContainer from "../components/ChartContainer";

function Requests({requestsEvents}) {
  return (
    <ChartContainer
      title={"Requests Cache Timeline"}
      description={
        "Tracking the requests arrive in since the server started"
      }
      data={requestsEvents}
      bgColor={"#9966cc"}
    />
  );
}

export default Requests;
