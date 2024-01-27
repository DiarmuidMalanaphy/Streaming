# Video Streaming Client-Server Module

This repository contains a client-server model designed for video streaming. It captures video data from a camera and transmits it to a server, which then relays the information to any client requesting the video feed.

The main conceit of this project is to develop local area video streaming for an upcoming project, again I've used Golang because it's performant, designed for concurrent systems, and I like working with it.

['API Description'][#API Documentation#]
['Installation and running the code][#Installation and Setup]



## Key Features

- **Client-Server Architecture**: Efficiently handles video data transmission from a camera to a server and then to clients.
- **Global Functionality**: The system is built to operate globally, with varying ping times depending on location, be aware you will have to handle portforwarding for it to work globally.
- **Multi-language Development**: The server is implemented in Go (Golang), and the client is in Python, with a future plan to migrate the client to Dart for integration with Flutter for a better user interface.
- **Socket Programming**: Uses low-level socket interfaces for precise control over network communications.

## Project Structure

The project is divided into several key components:

- **Server (Go)**: Handles incoming video streams, processes requests, and manages data distribution to clients.
- **Client (Python/Dart)**: Connects to the server, requests video data, and displays the video stream.
- **Networking**: Utilizes UDP for transmitting video data and handling client-server communications.



# API Documentation

## Overview
This API facilitates interactions with a video streaming server, enabling camera initialization, removal, feed requests, and updates. It is optimized for efficient video data transmission within a client-server architecture.

## API Calls

### `RequestTypeInitialiseCamera` - Request Type 1
- **Purpose**: Adds a new camera to the server's client list.
- **Input**: [`ImportedCamera`](#importedcamera)
  - Contains camera identifier, color bands, width, and height.
- **Output**: [`ExportedCamera`](#exportedcamera)
  - Returns the initialized camera details including name, ID, bands, width, and height.

### `RequestTypeRemoveCamera` - Request Type 2
- **Purpose**: Removes a camera from the server's client list.
- **Input**: Camera `ID`
  - Unique identifier of the camera to be removed.
- **Output**: None
  - Confirms the camera has been successfully removed.

### `RequestTypeUpdateCamera` - Request Type 3
- **Purpose**: Updates the server buffer with a new image from a camera.
- **Input**: [`IncomingImagePacket`](#incomingimagepacket)
  - Contains the camera ID and image data to be updated.
- **Output**: None
  - Indicates the image has been successfully updated in the buffer.

### `RequestTypeRequestFeed` - Request Type 4
- **Purpose**: Requests the video feed from a specific camera.
- **Input**: [`FeedRequest`](#feedrequest)
  - Contains the camera ID and sequence number for the requested feed.
- **Output**: [`FeedResponse`](#feedresponse)
  - Provides the most recent sequence number and buffer containing the feed data in UDP packets.




## Data Structures

### `ImportedCamera`
- Used when initialising the camera.
- **Name**: `[20]byte` - Identifier for the camera.
- **Bands**: `uint16` - Number of color bands.
- **Width**: `uint16` - Width resolution of the camera.
- **Height**: `uint16` - Height resolution of the camera.

### `FeedRequest`
- Used when requesting the feed from a single camera.
- **ID**: `uint16` - Unique identifier for the camera feed request.
- **SeqNum**: `uint32` - Sequence number for the feed request.

### `FeedResponse`
- Response to a feed request.
- **mostRecentSequenceNumber**: `uint32` - The most recent sequence number in the feed.
- **buffer**: `[][]UDPPacket` - Buffer holding the feed data in UDP packets.

### `ExportedCamera`
- Used when a request for camera information is requested.
- **Name**: `[20]byte` - Name of the camera.
- **ID**: `uint16` - Unique identifier for the camera.
- **Bands**: `uint16` - Number of color bands.
- **Width**: `uint16` - Width resolution of the camera.
- **Height**: `uint16` - Height resolution of the camera.

### `UDPPacket`
- Structure of a packet of video information. Sequence numbers and packet numbers are included as reassembly is required.
- **PacketNum**: `uint16` - Packet number within a sequence.
- **TotalPackets**: `uint16` - Total number of packets in a sequence.
- **SeqNum**: `uint32` - Sequence number of the packet.
- **Data**: `[packetSize]byte` - Data contained in the packet.

### `IncomingImagePacket`
- Used when sending information from the client to the server. Typically when a camera sends image information to the server.
- **CameraID**: `uint16` - Unique identifier for the camera sending the packet.
- **ImageInformation**: `ImagePacket` - Image data packet information.

### `ImagePacket`
- Used to describe an image
- **MessageID**: `uint32` - Unique identifier for the image message.
- **PacketNum**: `uint16` - Packet number within the image message.
- **TotalPackets**: `uint16` - Total number of packets in the image message.
- **SeqNum**: `uint32` - Number representing the image within the context of a larger video stream.
- **Data**: `[packetSize]byte` - Image data contained in the packet.


## Future Enhancements

- Transition the client application from Python to Dart for a more engaging UI using Flutter.
- Improve network efficiency and reduce latency for better global performance.



## Installation and Setup

1. Clone the repository to your local machine.
2. Ensure you have Go and Python installed, or Dart for the future client version.
3. Configure port forwarding on your network to allow external access to the server or just use the local network IP for local access.
4. Run the server using the provided batch file.
5. Connect clients by running the client script.

## Running the Code

To run the application, execute the provided batch files.




## Contributing

Contributions to the project are welcome! If you have ideas for improvements or want to report a bug, please open an issue or submit a pull request.

---

Feel free to adjust this README to better suit your project's specifics and to add any additional information that you think is relevant.
