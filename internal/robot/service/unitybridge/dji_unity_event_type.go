package unitybridge

type DJIUnityEventType int

const (
	SetValue DJIUnityEventType = iota
	GetValue
	GetAvailableValue
	PerformAction
	StartListening
	StopListening
	Activation
	LocalAlbum
	FirmwareUpgrade
	Connection         DJIUnityEventType = 100
	Security           DJIUnityEventType = 101
	PrintLog           DJIUnityEventType = 200
	StartVideo         DJIUnityEventType = 300
	StopVideo          DJIUnityEventType = 301
	Render             DJIUnityEventType = 302
	GetNativeTexture   DJIUnityEventType = 303
	VideoTransferSpeed DJIUnityEventType = 304
	AudioDataRecv      DJIUnityEventType = 305
	VideoDataRecv      DJIUnityEventType = 306
	NativeFunctions    DJIUnityEventType = 500
)

func DJIUnityEventTypes() []DJIUnityEventType {
	return []DJIUnityEventType{
		SetValue,
		GetValue,
		GetAvailableValue,
		PerformAction,
		StartListening,
		StopListening,
		Activation,
		LocalAlbum,
		FirmwareUpgrade,
		Connection,
		Security,
		PrintLog,
		StartVideo,
		StopVideo,
		Render,
		GetNativeTexture,
		VideoTransferSpeed,
		AudioDataRecv,
		VideoDataRecv,
		NativeFunctions,
	}
}
