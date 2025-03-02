## Echo
The message broker involves a pubsub server that handles messages sent by publishers via a REST API. When a publisher sends a message with a channel number and content, the server first saves the message to both cache and a file for persistence. This ensures that if the server unexpectedly shuts down, it can retrieve the messages from the file upon restart and deliver them to consumers once they reconnect. The server then broadcasts the message to any consumers currently listening on the specified channel. If no consumer is listening, the message is stored in the cache until a consumer connects. Upon receiving a message, a consumer sends an acknowledgment back to the server. Once acknowledged, the server deletes the message from both cache and file.
## Architecture Diagram
![Message broker](https://github.com/user-attachments/assets/58606f0a-8e38-4e05-bba2-d9cf1f77558e)

