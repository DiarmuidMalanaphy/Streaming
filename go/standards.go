package main

import "sync"

type ImportedCamera struct {
	Name   [20]byte
	Bands  uint16
	Width  uint16
	Height uint16
}

type FeedRequest struct {
	ID     uint16
	SeqNum uint32
}

type FeedResponse struct {
	mostRecentSequenceNumber uint32
	buffer                   [][]UDPPacket
}

type ExportedCamera struct {
	Name   [20]byte
	ID     uint16
	Bands  uint16
	Width  uint16
	Height uint16
}

func newExportedCamera(c Camera) ExportedCamera {
	// Initialize the image slice with the specified dimensions.

	return ExportedCamera{
		Name:   c.Name,
		Width:  c.Width,
		Bands:  c.Bands,
		Height: c.Height,
		ID:     c.ID,
	}
}

type UDPPacket struct {
	PacketNum    uint16
	TotalPackets uint16
	SeqNum       uint32
	Data         [512]byte
}

type ImagePacket struct {
	MessageID    uint32
	PacketNum    uint16
	TotalPackets uint16
	SeqNum       uint32
	Data         [512]byte
}

type IncomingImagePacket struct {
	CameraID         uint16
	ImageInformation ImagePacket
}

type ImageData struct {
	Packets      map[uint16][]byte
	Received     uint16
	TotalPackets uint16
	Mutex        sync.Mutex
}
