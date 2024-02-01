package main

import (
	"fmt"
)

// -> this standard will have to be relayed to the python too
const (
	// Camera Requests
	RequestTypeInitialiseCamera = uint8(1)
	RequestTypeRemoveCamera     = uint8(2)
	RequestTypeUpdateCamera     = uint8(3)

	RequestTypeRequestFeed    = uint8(4)
	RequestTypeRequestCameras = uint8(5)

	RequestSuccessful = uint8(200)
	RequestFailure    = uint8(255)
)

func main() {
	requestChannel := make(chan networkData)
	cameraMap := newCameraMap()

	go listen(requestChannel)

	for {
		select {
		case req := <-requestChannel:

			switch req.Request.Type {
			case RequestTypeInitialiseCamera:

				var ic ImportedCamera
				err := deserialiseData(req.Request.Payload, &ic)

				if err != nil {
					fmt.Println(err)
					outgoingReq, _ := generateRequest(ic, RequestFailure)
					sendUDP(req.Addr.String(), outgoingReq)
				}

				newCamera := cameraMap.addCamera(ic.Name, ic.Bands, ic.Width, ic.Height)
				outgoingReq, err := generateRequest(newCamera, RequestSuccessful)
				if err != nil {
					fmt.Println(err)
				}
				sendUDP(req.Addr.String(), outgoingReq)

			case RequestTypeRemoveCamera:
				var c Camera
				err := deserialiseData(req.Request.Payload, &c)

				if err != nil {
					fmt.Println(err)
					outgoingReq, _ := generateRequest(c, RequestFailure)
					sendUDP(req.Addr.String(), outgoingReq)
				}

				cameraMap.removeCamera(c)

				outgoingReq, _ := generateRequest(c, RequestSuccessful)
				sendUDP(req.Addr.String(), outgoingReq)

			case RequestTypeUpdateCamera:

				var i IncomingImagePacket

				err := deserialiseData(req.Request.Payload, &i)

				if err != nil {
					fmt.Println(err)
					outgoingReq, _ := generateRequest(i, RequestFailure)
					sendUDP(req.Addr.String(), outgoingReq)
					break
				}
				cameraID := i.CameraID
				camera, exists := cameraMap.getCamera(cameraID)
				if exists {
					go camera.handleIncomingPacket(i.ImageInformation)

				} else {
					fmt.Println("Camera not found")

				}

			case RequestTypeRequestFeed:
				var incoming FeedRequest
				err := deserialiseData(req.Request.Payload, &incoming)

				if err != nil {

					outgoingReq, _ := generateRequest(incoming, RequestFailure)
					sendUDP(req.Addr.String(), outgoingReq)
				}
				if cameraMap != nil {
					cameraID := incoming.ID
					camera, exists := cameraMap.getCamera(cameraID)
					if exists {

						packetLists := camera.getFeed(incoming.SeqNum)

						for _, packetList := range packetLists {
							for _, packet := range packetList {
								outgoingReq, err := generateRequest(packet, RequestSuccessful)
								if err != nil {
									fmt.Println(err)
									continue
								}

								err = sendUDP(req.Addr.String(), outgoingReq)
								if err != nil {
									fmt.Println(err)
								}
							}
						}
					} else {
						fmt.Println("Camera not found")
					}

				}

			case RequestTypeRequestCameras:
				if cameraMap != nil {

					cameras := cameraMap.getCameras()
					seqnum := 0
					for _, camera := range cameras {
						seqNum := uint32(seqnum + 1) // Generate a unique sequence number for this camera
						data := SerializeCamera(camera)
						packets := SplitIntoUDPPackets(seqNum, data)
						for _, packet := range packets {
							outgoingReq, err := generateRequest(packet, RequestSuccessful)
							err = sendUDP(req.Addr.String(), outgoingReq)
							if err != nil {
								fmt.Println("Error sending UDP packets:", err)
							}
						}

					}

				}

			}

		}
	}

}
