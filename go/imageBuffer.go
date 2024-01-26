package main

import "sort"

type ImageBufferItem struct {
	SeqNum  uint32
	Packets []UDPPacket
}

type ImageBuffer struct {
	Items    []ImageBufferItem
	Capacity int
}

func NewImageBuffer(capacity int) *ImageBuffer {
	return &ImageBuffer{
		Items:    []ImageBufferItem{},
		Capacity: capacity,
	}
}

func (b *ImageBuffer) Add(item ImageBufferItem) {
	// Add the new item to the buffer
	b.Items = append(b.Items, item)

	// Sort the buffer by sequence numbers in descending order
	sort.Slice(b.Items, func(i, j int) bool {
		return b.Items[i].SeqNum > b.Items[j].SeqNum
	})

	// Ensure the buffer size does not exceed the capacity
	// Keep only the items with the largest sequence numbers
	if len(b.Items) > b.Capacity {
		b.Items = b.Items[:b.Capacity]
	}
}
func (b *ImageBuffer) getPackets(seqNum uint32) [][]UDPPacket {
	var allPackets [][]UDPPacket
	for _, item := range b.Items {
		if item.SeqNum >= seqNum {
			allPackets = append(allPackets, item.Packets)
		}
	}
	return allPackets
}

func (b *ImageBuffer) GetSmallestSeqNum() *ImageBufferItem {
	if len(b.Items) == 0 {
		return nil
	}
	smallest := b.Items[0]
	for _, item := range b.Items {
		if item.SeqNum < smallest.SeqNum {
			smallest = item
		}
	}
	return &smallest
}
