import io
import math
import os
import socket
import struct
import time
from PIL import Image
import numpy as np
from PIL import Image
import numpy as np
import matplotlib.pyplot as plt
import cv2

from requestType import RequestType
from standardFormats import StandardFormats


class Networking:
    def __init__(self,hostIP,hostPort = 8000):
        self.hostIP = hostIP
        self.hostPort = hostPort


    def update_camera(self, img, camera_id,seq_num):
        #Essentially just break the image down into chunks and push them to the server
        packets = self.serialise_image(img, camera_id,seq_num)
        for packet in packets:
            #You can't have the socket_length lower than 0.01 otherwise it complains about blocking <- ideally you want it quite low though
            response = self.__send_general_payload_request(RequestType.RequestTypeUpdateCamera.value, packet,socket_length= 0.01)
            

    def initialise_camera(self,name,bands,width,height):
        #The name has to be a fixed size for transmission to the server so I've allocated it 20 characters, the excess will be removed.
        name_bytes = name.encode('utf-8')[:20]
        name_bytes += b'\x00' * (20 - len(name_bytes))
        #We create a payload for the initiation data
        payload = struct.pack("20sHHH", name_bytes, bands, width, height)
        response = self.__send_general_payload_request(RequestType.RequestTypeInitialiseCamera.value,payload)
        #I haven't done serious type checking here, this should be migrated to flutter or kotlin.
        type, _ , payload = self.deserialise_request(response)
        #We expect to get a receipt of the name, the ID associated and the dimensions of the image.
        #There is a maximum size on the server-side.
        unpacked_data = struct.unpack('20sHHHH', payload)
        return(unpacked_data)
        

    def request_feed(self,ID,seq_num):
        #Prepare the ID
        
        try:
            _,payload = self.serialize_payload([ID,seq_num],"HI")

            payload,new_seq = self._send_multiple_packet_request(RequestType.RequestTypeRequestFeed.value,payload)
            if payload is not None:
                images = []
                for seq_num, image_data in payload.items():
                    try:
                        # Deserialize the image using the original dimensions
                        #image = self.deserialise_image(image_data, bands, width, height)
                        images.append(image_data)  # Add the deserialized image to the list
                    except Exception as e:
                        print(f"Error deserializing image for sequence number {seq_num}: {e}")
                        return(None)
                return(images,seq_num)
        except Exception as e:
            return(None)

        
    


    
    def __send_general_payload_request(self,request_type,payload,socket_length = 0.5):
        # for more extensive documentation read __send_player_request
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
            try:
                
                # Payload length (int32)
                # We have to put self in on the python side as we need it for serialisation on go-side
                # Holy shit this was a pain in the ass, spent ages on wireshark figuring out that
                # 4 bytes were missing and what they were
                # Was the payload length 

                # It's to do with the way binarisation works in go, you have to define a fixed size because binarisation does not like unknown sized onjects
                # for generalisation i had to make the payload an undefined size in go, but i have to keep a standard.
                payload_length = len(payload)
                request_data = struct.pack(StandardFormats.RequestHeader.value, request_type, payload_length) + payload
                # 'B' for request type (uint8) and 'I' for payload length (int32)
                # The byte length should be 19
                 # Player length should be 14 + 1(type) + payload length (4)
                

                sock.settimeout(socket_length)
                sock.sendto(request_data, (self.hostIP, self.hostPort))
                
                try:
                    response = sock.recv(2048)
                    
                    
                    
                    return(response)
                    
                except socket.timeout:
                    # print("Timeout: No response received")
                    return((None,None))
                

            except socket.error as e:
                print(f"Socket error: {e}")
            except Exception as e:
                print(f"Other exception: {e}")


    def _send_multiple_packet_request(self, request_type, payload):
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as sock:
            try:
                request_data = struct.pack(StandardFormats.RequestHeader.value, request_type, len(payload)) + payload
                sock.settimeout(0.3)
                sock.sendto(request_data, (self.hostIP, self.hostPort))

                packets_by_seq_num = {}
                largest_seq_num = 0  # Initialize largest sequence number
                while True:
                    try:
                        response, _ = sock.recvfrom(2048)
                        _, _, payload = self.deserialise_request(response)
                        header = payload[:8]
                        packet_num, total_packets, seq_num = struct.unpack('=HHI', header)
                        largest_seq_num = max(largest_seq_num, seq_num)  # Update largest sequence number
                        if seq_num not in packets_by_seq_num:
                            packets_by_seq_num[seq_num] = {}
                        packets_by_seq_num[seq_num][packet_num] = payload[8:]
                        if all(len(packets_by_seq_num[sn]) == total_packets for sn in packets_by_seq_num):
                            break
                    except socket.timeout:
                        print("Timeout: No response received")
                        break

                images = {}
                for seq_num, packets in packets_by_seq_num.items():
                    if len(packets) == total_packets:
                        full_data = b''.join(packets[i] for i in range(total_packets))
                        images[seq_num] = full_data
                    else:
                        print(f"Error: Missing some packets for sequence number {seq_num}")

                return images, largest_seq_num

            except socket.error as e:
                print(f"Socket error: {e}")
                return None, None
            except Exception as e:
                print(f"Other exception: {e}")
                return None, None

    
    

    def serialize_payload(self,payload, format_string):
        # Ensure the format string starts with '='
        
        if not format_string.startswith('='):
            format_string = '=' + format_string

        try:
            # Convert payload to tuple if it has a to_tuple method
            if hasattr(payload, 'to_tuple'):
                payload = payload.to_tuple()

            # Serialize the payload
            serialized_data = struct.pack(format_string, *payload)

            # Calculate and return payload length and serialized data
            payload_length = len(serialized_data)
            return payload_length, serialized_data

        except struct.error as e:
            
            raise ValueError(f"Payload does not match format '{format_string}': {e}")
        

   
        


    def deserialize_payload(self,serialized_data, single_item_format): # <- works for several of the same
        # Ensure the format string starts with '='
        if not single_item_format.startswith('='):
            single_item_format = '=' + single_item_format

        # Calculate the size of a single item
        item_size = struct.calcsize(single_item_format)
        
        # Initialize an empty list to store unpacked items
        items = []

        # Iterate over the serialized data and unpack each item
        for i in range(0, len(serialized_data), item_size):
            # Extract a chunk of data for a single item
            item_data = serialized_data[i:i + item_size]

            try:
                # Unpack the item and append to the list
                unpacked_item = struct.unpack(single_item_format, item_data)
                items.append(unpacked_item)
            except struct.error as e:
                raise ValueError(f"Error unpacking data: {e}")

        return items
    


    def deserialise_request(self, request_data):
        # Unpack the first 5 bytes for Type and payloadLength
        if request_data is None:
            return(None)
        
        type, payload_length = struct.unpack(StandardFormats.RequestHeader.value, request_data[:5])

        # Extract the payload using the payload_length
        payload = request_data[5:5+payload_length]

        return (type, payload_length, payload)

    def serialise_image(self, img, camera_id,seq_num):
        # Convert the image to a numpy array and flatten it

        img_array = np.array(img)
        
        _, compressed_img = cv2.imencode('.jpg', img_array, [cv2.IMWRITE_JPEG_QUALITY, 50])
        img_bytes = compressed_img.tobytes()

        # Determine the total number of packets needed
        packet_size = 1312
        total_packets = math.ceil(len(img_bytes) / packet_size)

        # Generate a unique message ID
        
        
        message_id = os.urandom(4)  # Random 4-byte message ID
        message_id = struct.unpack('I', message_id)[0]#
        
        packets = []
        for i in range(total_packets):
            start = i * packet_size
            end = start + packet_size
            data = img_bytes[start:end]
            
            # Create an ImagePacket -> the standard is on the go-side
            #I had an issue with this. The MTU for a lot of networks is around 600kb 
            packet = struct.pack('=1312s',data)
            
            #My network was chopping off the packets after 1024kb and I couldn't figure out why.
            
            incoming_image_packet = struct.pack('=H', camera_id) + struct.pack("=IHHI", message_id,i,total_packets,seq_num) + packet 
            packets.append(incoming_image_packet)
        
        return(packets)
    
    def deserialise_image(self, payload, bands, width, height):
        payload = struct.unpack('<' + 'B' * len(payload), payload)
        image_data = payload
        image = [[[0 for _ in range(bands)] for _ in range(width)] for _ in range(height)]

        i = 0
        for y in range(height):  # Iterate over height first
            for x in range(width):  # Then iterate over width
                for b in range(bands):  # Finally iterate over bands
                    if i < len(image_data):
                        image[y][x][b] = image_data[i]
                        i += 1

        return image




if __name__ == "__main__":
    
    
    img = Image.open(r"C:\Users\liver\Streaming\test_image.bmp").convert("RGB")#.convert('L')
    
    networkingTool = Networking("127.0.0.1")
    
    networkingTool.initialise_camera("Benni",3,250,250)

    networkingTool.update_camera(img,1)
    

    networkingTool.request_feed([1],3,250,250)

    

