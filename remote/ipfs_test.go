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
	fmt.Println("IPFS gateway reachable: running actual tests")
	data := []byte("12345")
	rpi, err := ipfs.Add(data)
	if err!=nil{
		t.Errorf("FAIL: could not add data: %v", err)
	}
	err = ipfs.Remove(rpi)
	if err != nil {
		t.Errorf("FAIL: could not unpin cid: %v", err)
	}
}
