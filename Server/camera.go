package main

import (
	"sync"
)

var globalPacketStore = make(map[uint32]*ImageData)
var globalStoreMutex = sync.Mutex{}

type Camera struct {
	ID   uint16
	Name [20]byte

	Packets []UDPPacket
	Bands   uint16
	Width   uint16
	Height  uint16
	Buffer  *ImageBuffer
}

func newCamera(name [20]byte, bands uint16, width uint16, height uint16, ID uint16) Camera {
	// Initialize the image slice with the specified dimensions.

	return Camera{
		ID:   ID,
		Name: name,

		Bands:  bands,
		Width:  width,
		Height: height,
		Buffer: NewImageBuffer(5),
	}
}

func (c *Camera) exportCamera() ExportedCamera {
	return (newExportedCamera(*c))
}

func (c *Camera) getUDPPackets() []UDPPacket {
	return (c.Packets)
}

func (c *Camera) handleIncomingPacket(packet ImagePacket) {
	messageID := packet.MessageID
	packetNum := packet.PacketNum
	totalPackets := packet.TotalPackets
	seqNum := packet.SeqNum
	data := packet.Data[:]
	//Checks if the sequence number is relevant to us (discard old packets)
	if !c.Buffer.IsValidPacket(seqNum) {
		return
	}

	// We should probably limit concurrent access :(
	// Issue is that it makes it quite slow.
	globalStoreMutex.Lock()
	imageData, exists := globalPacketStore[messageID]
	if !exists {
		// Initialize if it's a new message
		imageData = &ImageData{
			Packets:      make(map[uint16][]byte),
			Received:     0,
			TotalPackets: totalPackets}
		globalPacketStore[messageID] = imageData
	}

	// Store the packet data

	imageData.Packets[packetNum] = data
	imageData.Received++

	globalPacketStore[messageID] = imageData
	globalStoreMutex.Unlock()

	// Check if all packets are received
	if imageData.Received == totalPackets {

		// Reassemble the image
		var fullImage []byte
		for p := uint16(0); p < totalPackets; p++ {
			fullImage = append(fullImage, imageData.Packets[p]...)
		}

		// Convert the byte slice to the image format
		c.updateImageFromBytes(fullImage, seqNum)

		delete(globalPacketStore, messageID)
	}

}

func (c *Camera) updateImageFromBytes(jpegBytes []byte, seqNum uint32) {

	packets := c.SplitIntoUDPPackets(seqNum, jpegBytes)

	newItem := ImageBufferItem{
		SeqNum:  seqNum,
		Packets: packets,
	}
	c.Buffer.Add(newItem)
}

func (c *Camera) SplitIntoUDPPackets(seqNum uint32, jpegBytes []byte) []UDPPacket {
	// We flatten the image for transmission

	totalPackets := (len(jpegBytes) + packetSize - 1) / packetSize

	var packets []UDPPacket
	for i := 0; i < len(jpegBytes); i += packetSize {
		// Adjust this to modify the packet_length
		var packetData [packetSize]byte

		end := i + packetSize
		if end > len(jpegBytes) {
			end = len(jpegBytes)
		}

		// Set PacketNum and TotalPackets
		packetNum := uint16(i / packetSize)

		// Copy image data into the packet
		copy(packetData[:], jpegBytes[i:end])

		packets = append(packets, UDPPacket{
			PacketNum:    packetNum,
			TotalPackets: uint16(totalPackets),
			SeqNum:       seqNum,
			Data:         packetData,
		})
	}

	return packets
}

func (c *Camera) getFeed(seqNum uint32) [][]UDPPacket { // make this return based on the sequenceNumber

	return (c.Buffer.getPackets(seqNum))
}
