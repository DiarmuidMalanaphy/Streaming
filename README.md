# Video Streaming Client-Server Module

This repository contains a client-server model designed for video streaming. A camera captures video data and transmits it to a server, which then relays the information to any client requesting the video feed.

## Key Features

- **Client-Server Architecture**: Efficient transmission of video data from a camera to a server, and from the server to clients.
- **Global Functionality**: The system is designed to work globally, though it should be noted that ping times may vary based on the user's location.
- **Network Configuration**: Requires port forwarding for optimal operation.
- **Language and Technology**: The server is implemented in Go (Golang) for its effectiveness in handling network operations. The client is currently written in Python, with plans to migrate to Dart and integrate with Flutter for an improved user interface.
- **Socket Programming**: Utilizes low-level socket interfaces for direct control over network communication.

## Running the Code

To run the application, simply execute the provided batch files.
