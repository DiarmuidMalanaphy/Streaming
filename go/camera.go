package main

import (
	"bytes"
	"encoding/binary"
	"image/jpeg"
	"time"
)

var globalPacketStore = make(map[uint32]ImageData)

type Camera struct {
	ID                      uint16
	Name                    [20]byte
	Image                   [][][]uint8
	Packets                 []UDPPacket
	Bands                   uint16
	Width                   uint16
	Height                  uint16
	CurrentDisplayTimestamp int64 // Useful for detecting old packets
}

func newCamera(name [20]byte, bands uint16, width uint16, height uint16, ID uint16) Camera {
	// Initialize the image slice with the specified dimensions.
	image := make([][][]byte, height)
	for y := range image {
		image[y] = make([][]byte, width)
		for x := range image[y] {
			image[y][x] = make([]byte, bands)
		}
	}

	return Camera{
		ID:                      ID,
		Name:                    name,
		Image:                   image,
		Bands:                   bands,
		Width:                   width,
		Height:                  height,
		CurrentDisplayTimestamp: time.Now().Unix(),
	}
}

func (c *Camera) updateImage(image [][][]uint8) {
	c.Image = image
}

func (c *Camera) exportCamera() ExportedCamera {
	return (newExportedCamera(*c))
}

func (c *Camera) getUDPPackets() []UDPPacket {
	return (c.Packets)
}

func (c Camera) SplitIntoUDPPackets() []UDPPacket {
	// We flatten the image for transmission
	var imageData []byte
	for y := range c.Image {
		for x := range c.Image[y] {
			for b := range c.Image[y][x] {
				imageData = append(imageData, c.Image[y][x][b])
			}
		}
	}

	packetSize := 512 - 4 // Adjust for the size of PacketNum and TotalPackets
	totalPackets := (len(imageData) + packetSize - 1) / packetSize

	var packets []UDPPacket
	for i := 0; i < len(imageData); i += packetSize {
		//Adjust this to modify the packet_length
		var packetData [512]byte

		end := i + packetSize
		if end > len(imageData) {
			end = len(imageData)
		}

		// Set PacketNum and TotalPackets
		packetNum := uint16(i / packetSize)
		binary.LittleEndian.PutUint16(packetData[0:2], packetNum)
		binary.LittleEndian.PutUint16(packetData[2:4], uint16(totalPackets))

		// Copy image data into the packet
		copy(packetData[4:], imageData[i:end])

		packets = append(packets, UDPPacket{
			PacketNum:    packetNum,
			TotalPackets: uint16(totalPackets),
			Data:         packetData,
		})
	}

	return packets
}

func (c *Camera) handleIncomingPacket(packet ImagePacket) {
	messageID := packet.MessageID
	packetNum := packet.PacketNum
	totalPackets := packet.TotalPackets
	data := packet.Data[:]

	// We should probably limit concurrent access :(
	// Issue is that it makes it quite slow.

	imageData, exists := globalPacketStore[messageID]
	if !exists {
		// Initialize if it's a new message
		imageData = ImageData{
			Packets:      make(map[uint16][]byte),
			Received:     0,
			TotalPackets: totalPackets,
			Timestamp:    time.Now().Unix(),
		}
	} else {
		// Check if the incoming packet is older than the current display
		if imageData.Timestamp < c.CurrentDisplayTimestamp {
			return // Ignore older packets
		}
	}

	// Store the packet data
	imageData.Packets[packetNum] = data
	imageData.Received++

	// Put the modified struct back into the map
	globalPacketStore[messageID] = imageData

	// Check if all packets are received
	if imageData.Received == totalPackets {
		c.CurrentDisplayTimestamp = imageData.Timestamp
		// Reassemble the image
		var fullImage []byte
		for p := uint16(0); p < totalPackets; p++ {
			fullImage = append(fullImage, imageData.Packets[p]...)
		}

		// Convert the byte slice to the image format
		c.updateImageFromBytes(fullImage)

		delete(globalPacketStore, messageID)
	}

}

func (c *Camera) updateImageFromBytes(jpegBytes []byte) {
	// We sent the image over the web as a jpg as it increased the speed a LOT
	// Convert the byte slice into an image
	imgReader := bytes.NewReader(jpegBytes)
	img, err := jpeg.Decode(imgReader)
	if err != nil {
		// Handle error (e.g., log or return it)
		return
	}

	// If the original format is RGB, you need to convert it
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Initialize the image slice with the camera's dimensions
	newImage := make([][][]uint8, height)
	for y := range newImage {
		newImage[y] = make([][]uint8, width)
		for x := range newImage[y] {
			newImage[y][x] = make([]uint8, 3) // Assuming RGB
		}
	}

	// Fill the newImage with pixel data from img
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Convert RGBA to uint8 (0-255 range)
			newImage[y][x][0] = uint8(r >> 8)
			newImage[y][x][1] = uint8(g >> 8)
			newImage[y][x][2] = uint8(b >> 8)

		}
	}

	// Update the camera's image
	c.Image = newImage
	c.Packets = c.SplitIntoUDPPackets()
}
