package video

type Handler interface {
	HandleVideo(img *RGB)
}
