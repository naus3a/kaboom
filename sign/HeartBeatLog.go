package sign

import(
	"sync"
	"encoding/json"
)

type HeartBeatLog struct{
	mutex sync.RWMutex
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
	l.mutex.Lock()
	if hasIt{
		if hb.Epoch > old.Epoch{
			l.lastHeartBeats[id] = *hb
		}
	}else{
		l.lastHeartBeats[id] = *hb
	}
	l.mutex.Unlock()
}

// GetLastHeartBeat gets you the latest heartbeat for an identity; the bool is false if thr id ia not found
func (l *HeartBeatLog) GetLastHeartBeat(id string)(HeartBeat, bool){
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	hb, b := l.lastHeartBeats[id]
	return hb, b
}

// GetNumIds returns the number of identities recorded
func (l *HeartBeatLog) GetNumIds()int{
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return len(l.lastHeartBeats)
}

func (l *HeartBeatLog) IsExpired(s *ArmoredShare, now int64) bool{
	if s== nil {
		return false
	}
	hb, hasIt := l.GetLastHeartBeat(s.AuthKey)
	if !hasIt{
		return false
	}
	return hb.IsExpired(int64(s.TTL), now)
}

func (l *HeartBeatLog) Serialize()([]byte, error){
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return json.Marshal(struct{
		LastHeartBeats map[string]HeartBeat `json:"lastHeartBeats"`
	}{
		LastHeartBeats: l.lastHeartBeats,
	})
}

func DeserializeHeartBeatLog(data []byte)(*HeartBeatLog, error){
	proxy := struct {
		LastHeartBeats map[string]HeartBeat `json:"lastHeartBeats"`
	}{}
	err := json.Unmarshal(data, &proxy)
	if err != nil{
		return nil, err
	}
	return &HeartBeatLog{
		lastHeartBeats: proxy.LastHeartBeats,
	}, nil
}
