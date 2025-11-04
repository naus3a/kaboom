package remote

import(
	"time"
	shell "github.com/ipfs/go-ipfs-api"
)

type IpfsRemoteController struct {
	sh *shell.Shell
}

func NewIpfsRemoteController(gateway string) *IpfsRemoteController{
	s := shell.NewShell(gateway)
	s.SetTimeout(2*time.Second)
	return &IpfsRemoteController{
		sh: s,
	}
}

func (ipfs *IpfsRemoteController)Ping()error{
	_, err := ipfs.sh.ID()
	return err
}
