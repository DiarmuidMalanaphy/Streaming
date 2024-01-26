from networking import Networking
import time
import cv2

class Camera:
    def __init__(self, ID, bands, width, height):
        self.ID = ID
        self.bands = bands
        self.width = width
        self.height = height

def capture_webcam_images(camera, networkTool):
    fps = 10
    frame_interval = 1.0 / fps  # Interval between frames

    # Initialize the webcam (0 is the default camera)
    cap = cv2.VideoCapture(0)

    # Check if the webcam is opened successfully
    if not cap.isOpened():
        print("Error: Could not open webcam.")
        return

    try:
        seq_num = 0
        while True:
            
            # Capture a single frame
            ret, frame = cap.read()

            # Check if the frame was captured successfully
            if not ret:
                print("Error: Could not read frame from webcam.")
                break

            # Process the frame -> 
                #OpenCv stores it as BGR
                #Resize the frame to the size of the "camera"
            resized_frame = cv2.resize(frame, (camera.width, camera.height))
            frame_rgb = cv2.cvtColor(resized_frame, cv2.COLOR_BGR2RGB)
            networkTool.update_camera(frame_rgb, camera.ID,seq_num)

            # Wait for the frame interval
            seq_num = seq_num + 1
            time.sleep(frame_interval)
    except KeyboardInterrupt:
        # Exit the loop when interrupted (e.g., by pressing Ctrl+C)
        print("Stopping webcam capture.")

    # Release the webcam
    cap.release()

# Prompt the user to enter the IP address
ip_address = input("Please enter the IP address to connect to: ")

#The camera creation is handled by the Server, ID is returned.
tool = Networking(ip_address)
camera = tool.initialise_camera("Benni", 3, 200, 200)
camera = Camera(camera[1], camera[2], camera[3], camera[4])
print("Camera ID is ",camera.ID)
capture_webcam_images(camera, tool)