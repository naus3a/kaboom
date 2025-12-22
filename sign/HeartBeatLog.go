package sign

type HeartBeatLog struct{
	lastHeartBeats []HeartBeat
}

func NewHeartBeatLog() *HeartBeatLog{
	return &HeartBeatLog{}
}
