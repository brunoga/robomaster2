package mobile

import "github.com/brunoga/robomaster2/modules/video"

type VideoHandler interface {
	HandleVideo(imgData []byte)
}

type Video struct {
	v *video.Video
}

func (v *Video) AddVideoHandler(videoHandler VideoHandler) (int, error) {
	return v.v.AddVideoHandler(&handler{videoHandler})
}

func (v *Video) StartSDCardRecording() {
	v.v.StartSDCardRecording()
}

func (v *Video) StopSDCardRecording() {
	v.v.StopSDCardRecording()
}

type handler struct {
	videoHandler VideoHandler
}

func (vh *handler) HandleVideo(img *video.RGB) {
	vh.videoHandler.HandleVideo(img.Pix)
}
