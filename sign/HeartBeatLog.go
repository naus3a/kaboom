package sign

type HeartBeatLog struct{
	lastHeartBeats map[string]HeartBeat
}

func NewHeartBeatLog() *HeartBeatLog{
	return &HeartBeatLog{
		lastHeartBeats: make(map[string]HeartBeat),
	}
}

// LogHeartBeat adds a new heartbeat to the log.
// if an heartbeat from the same author was already logged, it updates only if the new one is more rece.t.
// It is assumed that you already verified this heartbeat before passing it in.
func (l *HeartBeatLog) LogHeartBeat(id string, hb *HeartBeat){
	old, hasIt := l.GetLastHeartBeat(id)
	if hasIt{
		if hb.Epoch > old.Epoch{
			l.lastHeartBeats[id] = *hb
		}
	}else{
		l.lastHeartBeats[id] = *hb
	}
}

// GetLastHeartBeat gets you the latest heartbeat for an identity; the bool is false if thr id ia not found
func (l *HeartBeatLog) GetLastHeartBeat(id string)(HeartBeat, bool){
	hb, f := l.lastHeartBeats[id]
	return hb, f
}

// GetNumIds returns the number of identities recorded
func (l *HeartBeatLog) GetNumIds()int{
	return len(l.lastHeartBeats)
}
