from enum import Enum

class RequestType(Enum):
    RequestTypeInitialiseCamera = 1
    RequestTypeRemoveCamera = 2
    RequestTypeUpdateCamera = 3
    RequestTypeRequestFeed = 4
    RequestTypeRequestCameras = 5
    
    RequestSuccessful = 200
    RequestFailed = 255

