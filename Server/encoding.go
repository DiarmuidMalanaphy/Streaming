package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"io"
	"net"
	"reflect"
)

type networkData struct {
	Request Request
	Addr    net.Addr
}

func generateRequest(data interface{}, reqType uint8) ([]byte, error) {
	// First, serialize the data
	serializedData, err := serialiseData(data)
	if err != nil {
		return nil, err
	}

	// Create a Request with the serialized data as the payload
	req := newRequest(reqType, serializedData)
	serializedRequest, err := serialiseRequest(req)
	if err != nil {
		return nil, err
	}

	// Return the serialized request
	return serializedRequest, nil
}

func serialiseData(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		// Handle slice serialization
		for i := 0; i < v.Len(); i++ {
			err := binary.Write(buf, binary.LittleEndian, v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
	} else {
		// Handle non-slice serialization
		err := binary.Write(buf, binary.LittleEndian, data)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

//Essentially the same process as serialisation except you get weird behaviour due to the fact it's really difficult to tell
// where the end of the payload actually is.

func deserialiseData(data []byte, dataType interface{}) error {
	buf := bytes.NewReader(data)
	v := reflect.ValueOf(dataType)

	// Check if dataType is a pointer
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dataType must be a pointer")
	}

	v = v.Elem()

	if v.Kind() == reflect.Slice {
		// Handle slice deserialization
		sliceElementType := v.Type().Elem()
		for {
			elemPtr := reflect.New(sliceElementType)
			err := binary.Read(buf, binary.LittleEndian, elemPtr.Interface())
			if err == io.EOF {
				break // End of data
			}
			if err != nil {
				return err
			}
			v.Set(reflect.Append(v, elemPtr.Elem()))
		}
	} else {
		// Handle non-slice deserialization
		err := binary.Read(buf, binary.LittleEndian, dataType)
		if err != nil {
			return err
		}
	}

	return nil
}

// Type 1byte -> PayloadLength -> 4bytes -> PayloadBytes -> payload length to end.
func serialiseRequest(req Request) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write the Type field
	if err := binary.Write(buf, binary.LittleEndian, req.Type); err != nil {
		return nil, err
	}

	// Write the length of the Payload
	payloadLength := int32(len(req.Payload))
	if err := binary.Write(buf, binary.LittleEndian, payloadLength); err != nil {
		return nil, err
	}

	// Write the Payload bytes
	if _, err := buf.Write(req.Payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deserialiseRequest(data []byte) (Request, error) {
	var req Request
	buf := bytes.NewReader(data)

	// Read the Type field
	if err := binary.Read(buf, binary.LittleEndian, &req.Type); err != nil {
		return Request{}, err
	}

	// Read the length of the Payload
	var payloadLength int32
	if err := binary.Read(buf, binary.LittleEndian, &payloadLength); err != nil {
		return Request{}, err
	}

	// Read the Payload bytes
	req.Payload = make([]byte, payloadLength)
	if _, err := buf.Read(req.Payload); err != nil {
		return Request{}, err
	}

	return req, nil
}

func deserialisePicture(data []byte) (uint16, [][][]uint8, error) {
	// Read the header for width, height, and number of bands
	var ID, width, height, numBands uint32
	header := bytes.NewReader(data[:16]) // Header is now 12 bytes long

	if err := binary.Read(header, binary.LittleEndian, &ID); err != nil {
		return 0, nil, err
	}
	if err := binary.Read(header, binary.LittleEndian, &width); err != nil {
		return 1, nil, err
	}
	if err := binary.Read(header, binary.LittleEndian, &height); err != nil {
		return 2, nil, err
	}
	if err := binary.Read(header, binary.LittleEndian, &numBands); err != nil {
		return 3, nil, err
	}

	// Remaining data is the raw pixel data
	imgData := data[16:]

	// Create the band-dimensional array
	imgArray := make([][][]uint8, height)
	for i := range imgArray {
		imgArray[i] = make([][]uint8, width)
		for j := range imgArray[i] {
			imgArray[i][j] = make([]uint8, numBands)
		}
	}

	// Fill the array with pixel data
	index := 0

	for y := uint32(0); y < height; y++ {
		for x := uint32(0); x < width; x++ {
			for b := uint32(0); b < numBands; b++ {
				imgArray[y][x][b] = imgData[index]
				index++
			}
		}
	}

	return uint16(ID), imgArray, nil
}
