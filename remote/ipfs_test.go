package remote_test

import(
	"fmt"
	"testing"
	"github.com/naus3a/kaboom/remote"
)

func TestIpfs(t *testing.T){
	const gateway = "localhost:5001"
	ipfs := remote.NewIpfsRemoteController(gateway)
	err := ipfs.Ping()
	if err!=nil{
		fmt.Printf("Could not reach gateway: %v.\nStart a kubo daemon on localhost to run IPFS tests.", err)
		return
	}
}
