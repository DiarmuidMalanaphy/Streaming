# Video Streaming Client-Server Module

This repository contains a client-server model designed for video streaming. It captures video data from a camera and transmits it to a server, which then relays the information to any client requesting the video feed.

The main conceit of this project is to develop local area video streaming for an upcoming project, again I've used Golang because it's performant and I like working with it.

## Key Features

- **Client-Server Architecture**: Efficiently handles video data transmission from a camera to a server and then to clients.
- **Global Functionality**: The system is built to operate globally, with varying ping times depending on location, be aware you will have to handle portforwarding for it to work globally.
- **Multilingual Development**: The server is implemented in Go (Golang), and the client is in Python, with a future plan to migrate to Dart for integration with Flutter for a better user interface.
- **Socket Programming**: Uses low-level socket interfaces for precise control over network communications.

## Running the Code

To run the application, execute the provided batch files.

## Project Structure

The project is divided into several key components:

- **Server (Go)**: Handles incoming video streams, processes requests, and manages data distribution to clients.
- **Client (Python/Dart)**: Connects to the server, requests video data, and displays the video stream.
- **Networking**: Utilizes UDP for transmitting video data and handling client-server communications.

## Installation and Setup

1. Clone the repository to your local machine.
2. Ensure you have Go and Python installed, or Dart for the future client version.
3. Configure port forwarding on your network to allow external access to the server or just use the local network IP for local access.
4. Run the server using the provided batch file.
5. Connect clients by running the client script.

## Future Enhancements

- Transition the client application from Python to Dart for a more engaging UI using Flutter.
- Improve network efficiency and reduce latency for better global performance.

## Contributing

Contributions to the project are welcome! If you have ideas for improvements or want to report a bug, please open an issue or submit a pull request.

---

Feel free to adjust this README to better suit your project's specifics and to add any additional information that you think is relevant.
