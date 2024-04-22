package types

type Receiver struct {
	ReceivedMessage      []byte
	ReadableReceivedMsgs []Message
}

func (r *Receiver) NewReceivedMessage(msg []byte) {
	r.ReceivedMessage = msg
}
