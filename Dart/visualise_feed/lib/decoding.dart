import 'dart:typed_data';

class ExportedCamera {
  String name;
  int id;
  int bands;
  int width;
  int height;

  ExportedCamera({
    required this.name,
    required this.id,
    required this.bands,
    required this.width,
    required this.height,
  });

  static ExportedCamera fromBytes(Uint8List bytes) {
    // Assuming the name is encoded in ASCII
    String name = String.fromCharCodes(bytes.sublist(0, 20)).trim();
    int id = _getUint16(bytes, 20);
    int bands = _getUint16(bytes, 22);
    int width = _getUint16(bytes, 24);
    int height = _getUint16(bytes, 26);

    return ExportedCamera(name: name, id: id, bands: bands, width: width, height: height);
  }
}


Uint8List serializeFeedRequest(int id, int seqNum) {
    final byteData = ByteData(11); // 1 byte for uint8, 4 bytes for int32 (payload length), 2 bytes for uint16, and 4 bytes for uint32

    var requestNumber = 4;
    var payloadLength = 6; // Length of ID (2 bytes) + SeqNum (4 bytes)

    // Serialize requestNumber as uint8
    byteData.setUint8(0, requestNumber);

    // Serialize payload length as int32
    byteData.setUint32(1, payloadLength, Endian.little);

    // Serialize ID as uint16
    byteData.setUint16(5, id, Endian.little);

    // Serialize SeqNum as uint32
    byteData.setUint32(7, seqNum, Endian.little);

    return byteData.buffer.asUint8List();
  }

  Uint8List serializeCameraInfoRequest() {
    final byteData = ByteData(7); // 1 byte for uint8 (request type) + 4 bytes for int32 (payload length)

    var requestNumber = 5; // Assuming 5 is the request type for requesting camera info
    var payloadLength = 5; // Payload length is 0 since there are no additional parameters

    // Serialize requestNumber as uint8
    byteData.setUint8(0, requestNumber);

    // Serialize payload length as int32
    byteData.setUint32(1, payloadLength, Endian.little);

    return byteData.buffer.asUint8List();
  }

int _getUint16(Uint8List bytes, int start) {
  return bytes.buffer.asByteData().getUint16(start, Endian.little);
}