import ChartContainer from "../components/ChartContainer";

function Messages({messagesEvents}) {

  return (
      <ChartContainer
        title={"Messages Cache Timeline"}
        description={
          "Tracking the messages arrive in and are removed from the cache since the server started"
        }
        data={messagesEvents}
        bgColor={"#eeff00"}
      />
  );
}

export default Messages;
