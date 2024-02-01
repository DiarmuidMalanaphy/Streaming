package main

type CameraMap struct {
	cameraMap    map[uint16]*Camera
	nextCameraID uint16
}

func newCameraMap() *CameraMap {
	c := CameraMap{
		cameraMap:    make(map[uint16]*Camera),
		nextCameraID: 1,
	}
	return &c
}

func (cm *CameraMap) addCamera(name [20]byte, bands uint16, width uint16, height uint16) ExportedCamera {

	camera := newCamera(name, bands, width, height, cm.nextCameraID)
	cm.cameraMap[cm.nextCameraID] = &camera
	cm.nextCameraID = cm.nextCameraID + 1
	return (camera.exportCamera())

}

func (cm *CameraMap) removeCamera(c Camera) {
	delete(cm.cameraMap, c.ID)
}

func (cm *CameraMap) getCamera(ID uint16) (*Camera, bool) {
	camera, exists := cm.cameraMap[ID]
	return camera, exists
}

func (cm *CameraMap) getCameras() []ExportedCamera {
	var cameras []ExportedCamera
	for _, camera := range cm.cameraMap {
		cameras = append(cameras, camera.exportCamera())
	}
	return cameras
}
