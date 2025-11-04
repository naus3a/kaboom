package remote

import (
	"bytes"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"time"
)

const serviceName = "IPFS"

type IpfsRemoteController struct {
	sh *shell.Shell
}

func NewIpfsRemoteController(gateway string) *IpfsRemoteController {
	s := shell.NewShell(gateway)
	s.SetTimeout(2 * time.Second)
	return &IpfsRemoteController{
		sh: s,
	}
}

func (ipfs *IpfsRemoteController) Ping() error {
	_, err := ipfs.sh.ID()
	return err
}

func (ipfs *IpfsRemoteController) Add(data []byte) (*RemotePayloadId, error) {
	cid, err := ipfs.sh.Add(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &RemotePayloadId{
		Service: serviceName,
		Id:      cid,
	}, nil
}

func (ipfs *IpfsRemoteController) Remove(rpi *RemotePayloadId) error {
	if rpi == nil {
		return fmt.Errorf("nil RemotePayloadId")
	}
	if rpi.Service != serviceName {
		return fmt.Errorf("wrong service: expected %s, got %s", serviceName, rpi.Service)
	}
	err := ipfs.sh.Unpin(rpi.Id)
	if err != nil {
		return err
	}

	return nil
}
