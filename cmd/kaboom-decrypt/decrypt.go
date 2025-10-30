package main

import(
	"os"
	"fmt"
	"flag"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/fs"
	"github.com/naus3a/kaboom/payload"
	"github.com/naus3a/kaboom/sign"
)

const usage = `Usage:
	kaboom-decrypt [-s a.shab,b.shab] [-k key.keyb] [-p encrypted.xyz] [-o outputName]

	Options:
		-h, --help	this help screen
		-v, --version	prints version
		-s, shares	a list of csv share paths
		-k, --key		the key file
		-p, --payload		the encrypted file
		-o, --output		ourput file name; default is 'decrypted'; key will use the .keyb eztension and shares will use  .shab
`

func main(){
	var hFlag bool
	var vFlag bool
	var sFlag string
	var kFlag string
	var pFlag string
	var oFlag string

	hasAtLeast1Task := false
	var key *payload.ArmoredPayloadKey = nil
	var err error = nil

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&sFlag, "", "s", "shares")
	cmd.AddArg(&kFlag, "", "k", "key")
	cmd.AddArg(&pFlag, "", "p", "payload")
	cmd.AddArg(&oFlag, "decrypted", "o", "output")
	flag.Parse()

	if hFlag{
		flag.Usage()
		os.Exit(0)
	}

	if vFlag{
		fmt.Println(cmd.Version)
		os.Exit(0)
	}

	if sFlag!=""{
		hasAtLeast1Task = true
		key, err = combineShares(sFlag, oFlag)
		cmd.ReportErrorAndExit(err)
	}

	if kFlag!=""{
		keySer, err := fs.LoadFile(kFlag)
		cmd.ReportErrorAndExit(err)
		key, err = payload.Deserialize(keySer)
		cmd.ReportErrorAndExit(err)
	}
	
	if pFlag!=""{
		hasAtLeast1Task = true
		if key==nil{
			fmt.Println("you need a valid key to decrypt")
			os.Exit(1)
		}
		data, err := fs.LoadFile(pFlag)
		cmd.ReportErrorAndExit(err)
		plain, err := key.Decrypt(data)
		cmd.ReportErrorAndExit(err)
		fName := oFlag+cmd.ExtPlain
		err = fs.SaveFile(plain, fName)
		cmd.ReportErrorAndExit(err)
	}

	if !hasAtLeast1Task {
		fmt.Println("you need to specify at least one task")
		flag.Usage()
		os.Exit(1)
	}
}

func combineShares(arg string, outName string)(*payload.ArmoredPayloadKey, error){
	pthShares, err := cmd.UnpackCsvArg(&arg)
	if err!=nil {
		return nil, err
	}

	signedShares := make([]*sign.ArmoredShare, len(pthShares))
	for i:=0;i<len(pthShares);i++{
		jsonData, err := fs.LoadFile(pthShares[i])
		if err != nil {
			return nil, err
		}
		signedShares[i], err = sign.DeserializeShare(jsonData)
		if err != nil {
			return nil, err
		}
	}
	
	verified, err := verifyShares(signedShares)
	if err!=nil{
		return nil, err
	}
	if !verified{
		return nil, fmt.Errorf("bad share signature")
	}

	shares := make([][]byte, len(signedShares))
	for i:=0; i<len(signedShares); i++{

		data, err := signedShares[i].GetData()
		if err != nil {
			return nil, err
		}
		shares[i] = data
	}
	key, err := payload.CombineSharesInArmoredPayloadKey(shares)
	if err!=nil {
		return nil, err
	}
	keySer, err := key.Serialize()
	if err!=nil{
		return nil, err
	}
	fName := outName+cmd.ExtKey
	err = fs.SaveFile(keySer, fName)
	if err!=nil{
		return nil, err
	}
	return key, nil
}

func verifyShares(shares []*sign.ArmoredShare)(bool, error){
	if len(shares)<2{
		return true, nil
	}

	for i:=1; i<len(shares);i++{
		b, err := shares[0].VerifyShare(shares[i])
		if err != nil{
			return b, err
		}
		if !b{
			return false, err
		}
	}
	return true, nil
}
