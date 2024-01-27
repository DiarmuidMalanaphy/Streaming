from networking import Networking
import time
import cv2
import numpy as np

def convert_rgb_to_bgr(image):
    if image is not None:
        # Convert RGB to BGR
        return cv2.cvtColor(image, cv2.COLOR_RGB2BGR)
    return None

def continuously_request_images(network_tool, ID, interval):
    seq_num = 0
    
    cv2.namedWindow("Video Feed", cv2.WINDOW_NORMAL)
    while True:
        try:
            results = network_tool.request_feed(ID, seq_num)
            if results is not None:
                jpeg_images, new_seq_num = results
                print(new_seq_num)
                num_images = len(jpeg_images)

                # Calculate the time each image should be displayed
                display_time_per_image = interval / max(1, num_images)

                for jpeg_image in jpeg_images:
                    np_arr = np.frombuffer(jpeg_image, dtype=np.uint8)
                    image = cv2.imdecode(np_arr, cv2.IMREAD_COLOR)

                    if image is not None:
                        # Convert RGB to BGR if necessary
                        bgr_image = convert_rgb_to_bgr(image)
                        cv2.imshow("Video Feed", bgr_image)
                        time.sleep(display_time_per_image)  # Adjusted display time
                        if cv2.waitKey(1) & 0xFF == ord('q'):
                            break  # Exit the inner loop if 'q' is pressed

                seq_num = new_seq_num

            time.sleep(interval)  # Wait for the specified interval
            if cv2.waitKey(1) & 0xFF == ord('q'):
                break  # Exit the outer loop if 'q' is pressed

        except KeyboardInterrupt:
            break  # Exit the loop on keyboard interrupt

    cv2.destroyAllWindows()

# Prompt the user to enter the IP address
ip_address = input("Please enter the IP address to connect to: ")
camera_ID = int(input("Please enter the ID: "))

# Example usage
tool = Networking(ip_address)
continuously_request_images(tool, camera_ID, 0.04)