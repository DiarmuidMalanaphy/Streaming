package main

type Request struct {
	Type    uint8
	Payload []byte // Raw data, can be interpreted based on the request type
}

func newRequest(Type uint8, Payload []byte) Request {
	req := Request{
		Type,
		Payload,
	}
	return req
}
