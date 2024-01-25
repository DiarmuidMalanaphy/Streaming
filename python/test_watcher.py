from networking import Networking
import time
import cv2
import numpy as np

def convert_to_numpy(image):
    if image is not None:
        # Convert the image to a NumPy array for RGB
        height = len(image)
        width = len(image[0])
        np_image = np.zeros((height, width, 3), dtype=np.uint8)
        for y in range(height):
            for x in range(width):
                for c in range(3):
                    np_image[y][x][c] = image[y][x][c]
        return np_image
    return (None)

def continuously_request_images(network_tool, ID, bands, height, width, interval):
    cv2.namedWindow("Video Feed", cv2.WINDOW_NORMAL)
    while True:
        try:
            
            result = network_tool.request_feed(ID, bands, height, width)
            if result is not None:
                image,timeStamp = result
                image = convert_to_numpy(image)
                print(timeStamp)  
                
                if image is not None:
                    # Assuming image is a NumPy array in the correct shape and dtype
                    cv2.imshow("Video Feed", image)
            time.sleep(interval)  # Wait for the specified interval
            if cv2.waitKey(1) & 0xFF == ord('q'):
                break  # Exit the loop if 'q' is pressed

        except KeyboardInterrupt:
            break  # Exit the loop on keyboard interrupt
        

    cv2.destroyAllWindows()

# Prompt the user to enter the IP address
ip_address = input("Please enter the IP address to connect to: ")
camera_ID = int(input("Please enter the ID: "))

# Example usage
tool = Networking(ip_address)
continuously_request_images(tool, [camera_ID], 3, 200, 200, 0.25)