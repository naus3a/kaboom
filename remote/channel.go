package remote

import(
	"fmt"
	"time"
	"crypto/sha256"
	"encoding/base64"
)

func MakeChannelName(now time.Time,salt string)string{
	today := now.Truncate(24*time.Hour).Unix()
	s := fmt.Sprintf("%d%s", today, salt)
	
	hasher := sha256.New()
	hasher.Write([]byte(s))
	fullHash := hasher.Sum(nil)
	shortHash := fullHash[:8]
	return base64.RawURLEncoding.EncodeToString(shortHash)
}

func MakeChannelNameNow(salt string)string{
	return MakeChannelName(time.Now().UTC(), salt)
}
