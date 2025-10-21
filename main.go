package main

import (
	"github.com/naus3a/kaboom/payload"
	"fmt"
)

func main(){
	k,_:= payload.MakeSymmetricKey()
	fmt.Println(k)
}
