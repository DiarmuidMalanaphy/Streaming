import 'dart:io';
import 'dart:typed_data';
import 'dart:math' as math;

class UDPHandler {
  Future<List<dynamic>> sendRequest(
    String ipAddress, Uint8List feedRequest,
    {int port = 8000}) async {

    Map<int, Map<int, Uint8List>> packetsBySeqNum = {}; // Sequence number -> (Packet number -> Packet data)
    Map<int, int> totalPacketsBySeqNum = {}; // Sequence number -> Total number of packets
    int maxSeqNum = 0; // Initialize maximum sequence number

    RawDatagramSocket socket = await RawDatagramSocket.bind(InternetAddress.anyIPv4, 0);
    socket.send(feedRequest, InternetAddress(ipAddress), port);

    try {
      await for (RawSocketEvent event in socket.timeout(const Duration(milliseconds: 300))) {
        if (event == RawSocketEvent.read) {
          Datagram? datagram = socket.receive();
          if (datagram != null) {
            Uint8List data = datagram.data;
            var result = deserialiseRequest(data);
            if (result != null) {
              int type = result.item1;
              int payloadLength = result.item2;
              Uint8List payload = result.item3;

              // Extract packet info from payload
              int packetNum = _getUint16(payload, 0);
              int totalPackets = _getUint16(payload, 2);
              int seqNum = _getUint32(payload, 4);
              Uint8List packetData = payload.sublist(8);

              // Store packet data by sequence number and packet number
              packetsBySeqNum.putIfAbsent(seqNum, () => {});
              packetsBySeqNum[seqNum]![packetNum] = packetData;
              totalPacketsBySeqNum[seqNum] = totalPackets;
              maxSeqNum = math.max(maxSeqNum, seqNum); // Update maximum sequence number
            }
          }
        }
      }
    } catch (e) {
      print('Error receiving data: $e');
    } finally {
      socket.close();
    }

    // Reassemble packets into feed responses
    Map<int, Uint8List> feedResponses = {};
    for (var seqNum in packetsBySeqNum.keys) {
      var packets = packetsBySeqNum[seqNum]!;
      if (packets.length == totalPacketsBySeqNum[seqNum]) {
        // All packets are present for this sequence number
        List<int> allData = [];
        for (var packetNum in packets.keys.toList()..sort()) {
          allData.addAll(packets[packetNum]!);
        }
        feedResponses[seqNum] = Uint8List.fromList(allData);
      } else {
        print("Incomplete image data for sequence number $seqNum");
      }
    }

    return [feedResponses, maxSeqNum];
  }
  int _getUint16(Uint8List data, int offset) {
    return ByteData.sublistView(data, offset, offset + 2).getUint16(0, Endian.little);
  }

  int _getUint32(Uint8List data, int offset) {
    return ByteData.sublistView(data, offset, offset + 4).getUint32(0, Endian.little);
  }

  Tuple3<int, int, Uint8List>? deserialiseRequest(Uint8List request_data) {
    if (request_data.length < 5) {
      return null; // Not enough data to process
    }
    int type = request_data[0];
    int payloadLength = _getUint32(request_data, 1);
    
    Uint8List payload = request_data.sublist(5, 5 + payloadLength);
    return Tuple3(type, payloadLength, payload);
  }
}

class Tuple3<T1, T2, T3> {
  final T1 item1;
  final T2 item2;
  final T3 item3;

  Tuple3(this.item1, this.item2, this.item3);
}