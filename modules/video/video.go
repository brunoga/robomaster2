package video

import (
	"encoding/binary"
	"fmt"
	"image"
	"sync"

	"github.com/brunoga/robomaster2/internal/robot/service"
	"github.com/brunoga/robomaster2/internal/robot/service/dji"
	"github.com/brunoga/robomaster2/internal/robot/service/unitybridge"
	"github.com/brunoga/robomaster2/support"
)

type Video struct {
	logger *support.Logger

	m             sync.Mutex
	videoHandlers map[int]Handler
	img           *RGB
}

func New(logger *support.Logger) *Video {
	return &Video{
		logger,
		sync.Mutex{},
		make(map[int]Handler),
		NewRGB(image.Rect(0, 0, 1280, 720)),
	}
}

func (v *Video) Start() error {
	ub := unitybridge.DJIUnityBridgeInstance()

	ub.RegisterEventHandler(v, unitybridge.GetNativeTexture)
	ub.RegisterEventHandler(v, unitybridge.VideoTransferSpeed)
	ub.RegisterEventHandler(v, unitybridge.VideoDataRecv)

	return nil
}

func (v *Video) Stop() error {
	ub := unitybridge.DJIUnityBridgeInstance()

	ub.UnregisterEventHandler(v)

	return nil
}

func (v *Video) AddVideoHandler(videoHandler Handler) (int, error) {
	v.m.Lock()
	defer v.m.Unlock()

	id := len(v.videoHandlers)

	if id == 0 {
		unitybridge.DJIUnityBridgeInstance().SendEventWithoutDataOrTag(
			unitybridge.NewDJIUnityEventWithType(unitybridge.StartVideo))
	}

	v.videoHandlers[id] = videoHandler

	unitybridge.DJIUnityBridgeInstance().SendEventWithoutDataOrTag(
		unitybridge.NewDJIUnityEventWithType(unitybridge.GetNativeTexture))

	return id, nil
}

func (v *Video) RemoveVideoHandler(id int) error {
	v.m.Lock()
	defer v.m.Unlock()

	_, ok := v.videoHandlers[id]
	if !ok {
		return fmt.Errorf("invalid video handler id: %v", id)
	}

	delete(v.videoHandlers, id)

	if len(v.videoHandlers) == 0 {
		unitybridge.DJIUnityBridgeInstance().SendEventWithoutDataOrTag(
			unitybridge.NewDJIUnityEventWithType(unitybridge.StopVideo))
	}

	return nil
}

func (v *Video) OnEventCallback(event *unitybridge.DJIUnityEvent, data []byte, tag uint64) {
	switch event.Type() {
	case unitybridge.GetNativeTexture:
		v.logger.INFO("GetNativeTexture: %v", string(data))
		// TODO(bga): Set correct texture resolution.
		v.m.Lock()
		v.img = NewRGB(image.Rect(0, 0, 1280, 720))
		v.m.Unlock()
	case unitybridge.VideoTransferSpeed:
		value := binary.NativeEndian.Uint64(data)
		v.logger.INFO("VideoTransferSpeed: %v", value)
		// TODO(bga): What to do here?
	case unitybridge.VideoDataRecv:
		v.m.Lock()

		v.img.Pix = data

		for _, videoHandler := range v.videoHandlers {
			videoHandler.HandleVideo(v.img)
		}

		v.m.Unlock()
	}
}

func (v *Video) StartSDCardRecording() {
	cc := service.DJICommandControllerInstance()
	cc.GetValueForKey(dji.DJICameraMode, func(result *dji.DJIResult) {
		if !result.Succeeded() {
			// Could not get camera mode. Nothing else to do.
			v.logger.ERROR("Failed to get camera mode: %v", result)
			return
		}
		if result.Value().(dji.DJILongParamValue).Value != 1 {
			// Camera not in video more. Change it.
			cc.SetValueForKeyWithNumber(dji.DJICameraMode, 1, func(result *dji.DJIResult) {
				if !result.Succeeded() {
					// Could not set camera mode. Nothing else to do.
					v.logger.ERROR("Failed to set camera mode: %v", result)
					return
				}
				cc.PerformAction(dji.DJICameraStartRecordVideo, func(result *dji.DJIResult) {
					if !result.Succeeded() {
						// Could not start recording. Nothing else to do.
						v.logger.ERROR("Failed to start recording: %v", result)
						return
					}
				})
			})
		} else {
			// Camera already in video mode. Start recording.
			cc.PerformAction(dji.DJICameraStartRecordVideo, func(result *dji.DJIResult) {
				if !result.Succeeded() {
					v.logger.ERROR("Failed to start recording: %v", result)
					return
				}
			})
		}
	})
}

func (v *Video) StopSDCardRecording() {
	cc := service.DJICommandControllerInstance()
	cc.PerformAction(dji.DJICameraStopRecordVideo, func(result *dji.DJIResult) {
		if !result.Succeeded() {
			v.logger.ERROR("Failed to stop recording: %v", result)
			return
		}
	})
}