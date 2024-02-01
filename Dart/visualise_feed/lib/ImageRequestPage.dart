import 'dart:async';
import 'dart:ffi';
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:visualise_feed/AnimatedRays.dart';
import 'package:visualise_feed/UDPHandler.dart';
import 'package:visualise_feed/decoding.dart';

class ImageRequestPage extends StatefulWidget {
  final String ipAddress;

  ImageRequestPage({Key? key, required this.ipAddress}) : super(key: key);

  @override
  _ImageRequestPageState createState() => _ImageRequestPageState();
}

class _ImageRequestPageState extends State<ImageRequestPage> {
  
  final _controller = TextEditingController();
  String _serializedData = '';
  Image? _receivedImage; // Image widget to display the received JPEG image
  int _seqNum = 0; // Sequence number
  bool _isRequesting = false; 
  final List<Uint8List> _imageBuffer = [];
  Timer? _updateTimer;
  ImageProvider? _currentImageProvider;
  List<int> _sortedCameraIds = [];
  int _previousCameraId = -1;
  bool _isLoading = false;

Future<bool> requestCameraInfo() async {
  Uint8List serializedData = serializeCameraInfoRequest();
  UDPHandler udpHandler = UDPHandler();

  List<dynamic> cameraInfo = await udpHandler.sendRequest(widget.ipAddress, serializedData);
  Map<int, Uint8List> cameraInfoResponses = cameraInfo[0];
  List<ExportedCamera> cameras = [];
  for (var responseBytes in cameraInfoResponses.values) {
    ExportedCamera camera = ExportedCamera.fromBytes(responseBytes);
    cameras.add(camera);
  }

  if (cameras.isNotEmpty) {
    // Sort cameras by ID and update the sorted list
    _sortedCameraIds = cameras.map((camera) => camera.id).toList()..sort();
    return true;
  }
  return false;
}

  @override
  void initState() {
    super.initState();
    _updateTimer = Timer.periodic(Duration(milliseconds: 100), (timer) {
      // Check if there are images in the buffer
      if (_imageBuffer.isNotEmpty) {
        // Display the first image in the buffer
        _updateImage(_imageBuffer.first);
        // Remove the displayed image from the buffer
        _imageBuffer.removeAt(0);
      } else {
        // Optionally, stop the timer if there are no images left
        // timer.cancel();
      }
    });
  }


Future<int> getNextCameraId(int currentId) async {
  bool isLoaded = await requestCameraInfo();
  if (isLoaded && _sortedCameraIds.isNotEmpty) {
    int index = _sortedCameraIds.indexOf(currentId);
    if (index == -1 || index == _sortedCameraIds.length - 1) {
      // If currentId is not found or is the last one, return the first ID
      return _sortedCameraIds.first;
    } else {
      // Return the next ID in the list
      return _sortedCameraIds[index + 1];
    }
  }
  return currentId; // Return the current ID if unable to get camera info or if the list is empty
}

  @override
  void dispose() {
    _updateTimer?.cancel(); // Cancel the timer if it's running
    super.dispose();
  }

// Tracks if we are continuously requesting images

  
  void continuouslyRequestImages(int id) async {
    _isRequesting = true;
    UDPHandler udpHandler = UDPHandler();

    while (_isRequesting) {
      Uint8List serializedData = serializeFeedRequest(id, _seqNum);
      List<dynamic> response = await udpHandler.sendRequest(widget.ipAddress, serializedData);
      Map<int, Uint8List> feedResponses = response[0];
      _seqNum = response[1];

      if (feedResponses.isNotEmpty) {
        // Sort the packets based on sequence numbers
        var sortedSeqNums = feedResponses.keys.toList()..sort();
        _imageBuffer.clear();
        for (var seqNum in sortedSeqNums) {
          _imageBuffer.add(feedResponses[seqNum]!);
        }
      }

      await Future.delayed(Duration(milliseconds: 20)); // Adjust as needed for server response time
    }
  }

  void toggleContinuousRequest() {
    
    if (_controller.text.isNotEmpty) {
      int id = int.tryParse(_controller.text) ?? 0;

      if (_isRequesting) {
        _isRequesting = false; // Stop requesting
      } else {
        continuouslyRequestImages(id); // Start requesting
      }
    }
  }
    Future<void> updateCameraIdAndStartStream(int newId) async {
    _previousCameraId = int.tryParse(_controller.text) ?? 0;

    if (_isRequesting) {
      toggleContinuousRequest(); // Stop the current stream
    }

    setState(() {
      _isLoading = true; // Start loading
      _currentImageProvider = MemoryImage(kTransparentImage); // Reset image to default
    });

    await Future.delayed(Duration(milliseconds: 100)); // Give some time for the server to process
    _controller.text = newId.toString(); // Update the controller text with new ID
    toggleContinuousRequest(); // Start the new stream
  }

  void _updateImage(Uint8List newImageBytes) {
    if (!mounted) return;
    setState(() {
      _currentImageProvider = MemoryImage(newImageBytes);
      _isLoading = false; // Stop loading
    });
  }

  @override
  
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('StreamBeam'),
        centerTitle: true,
        backgroundColor: Colors.deepPurple,
        elevation: 0,
      ),
      body: AnimatedBackground(
        
        child: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [Colors.deepPurple, Colors.purpleAccent],
            ),
          ),
          child: Padding(
            padding: EdgeInsets.symmetric(horizontal: 16.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: <Widget>[
              Padding(
                padding: EdgeInsets.all(20),
                child: Text(
                  'Enter Camera ID:',
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white),
                  textAlign: TextAlign.center,
                ),
              ),
              Padding(
                padding: EdgeInsets.symmetric(horizontal: 30),
                child: TextField(
                  controller: _controller,
                  onChanged: (newText) {
                    setState(() {
                      _imageBuffer.clear();
                      _currentImageProvider = null; // Set image to null when a new ID is typed
                    });
                  },
                  decoration: InputDecoration(
                    labelText: 'Camera ID',
                    hintText: 'e.g., 101',
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(10),
                    ),
                    prefixIcon: Icon(Icons.videocam),
                  ),
                  keyboardType: TextInputType.number,
                ),
              ),
              SizedBox(height: 20),
              ElevatedButton(
                onPressed: toggleContinuousRequest,
                child: Text(_isRequesting ? 'Stop Stream' : 'Start Stream'),
                style: ElevatedButton.styleFrom(
                  primary: Colors.deepPurple,
                  onPrimary: Colors.white,
                  padding: EdgeInsets.symmetric(horizontal: 30, vertical: 15),
                  textStyle: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                ),
              ),
              SizedBox(height: 10),
              Divider(
                color: Colors.white.withOpacity(0.5),
                thickness: 2,
                endIndent: 30,
                indent: 30,
              ),
              Expanded(
                flex: 10,
                child: GestureDetector(
                  onHorizontalDragEnd: (details) async {
                    if (details.primaryVelocity! < 0) {
                      // Swipe left: load the next camera
                      if (_controller.text != "") {
                        int nextId = await getNextCameraId(int.parse(_controller.text));
                        if (nextId != int.parse(_controller.text)) {
                          await updateCameraIdAndStartStream(nextId);
                        }
                      }
                    }
                  },
                  child: AnimatedSwitcher(
                      duration: Duration(milliseconds: 300),
                      transitionBuilder: (Widget child, Animation<double> animation) {
                        int currentId = int.tryParse(_controller.text) ?? 0;
                        bool isIdChanged = currentId != _previousCameraId;
                        
                        return isIdChanged ? SlideTransition(
                          position: Tween<Offset>(begin: Offset(1.0, 0.0), end: Offset(0.0, 0.0)).animate(animation),
                          child: child,
                        ) : child;
                      },
                    child: Container(
                      key: ValueKey<int>(int.tryParse(_controller.text) ?? 0),
                      margin: EdgeInsets.all(20), // Adjusted margin
                      decoration: BoxDecoration(
                        border: Border.all(color: Colors.white.withOpacity(0.5), width: 3),
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(10),
                        child: _isLoading
                            ? Center(child: CircularProgressIndicator())
                            : (_currentImageProvider != null 
                                ? FadeInImage(
                                    placeholder: MemoryImage(kTransparentImage),
                                    image: _currentImageProvider!,
                                    fit: BoxFit.cover,
                                  )
                                : Center(child: Text('No Image', style: TextStyle(color: Colors.white)))
                              ),
                      ),
                    ),
                  ),
                ),
              ),
              SizedBox(height: 10),
            ],
          ),
        ),
      ),
    ),
  );
}
}
final Uint8List kTransparentImage = Uint8List.fromList([0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82]);
