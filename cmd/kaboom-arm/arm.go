package main

import (
	"flag"
	"fmt"
	"github.com/naus3a/kaboom/cmd"
	"github.com/naus3a/kaboom/fs"
	"github.com/naus3a/kaboom/payload"
	"github.com/naus3a/kaboom/sign"
	"github.com/naus3a/kaboom/remote"
	"os"
)

const usage = `Usage:
kaboom-arm -p plaintext.file [-l localEncrypted.file] [-s 3] [-t 2] [-k signingkeys.file] [-g localhost:5001]

Options:
	-h, --help				this help screen
	-v, --version			prints version
	-p, --payload			the file you want to encrypt
	-l, --local			the local version of the encrypted payload output
	-g, --gateway			thr ipfs gateway
	-s, --shares			number of shares (default: 3)
	-t, --threshold		share threshold (default: 2)
	-n, --notes				extra notes for your payload
	-d, --delete			secure-delete plaintext
	-m, --maxTtl		maximum time in hours since last heartbeat (default: 24)
	-k, --signingkeys	the output file containing your signing keys (default: signingkeys.sigb)
`

func main() {
	var pFlag string
	var lFlag string
	var nFlag string
	var kFlag string
	var gFlag string
	var sFlag uint
	var mFlag uint
	var tFlag uint
	var vFlag bool
	var hFlag bool
	var dFlag bool

	cmd.InitCli(usage)
	cmd.AddArg(&hFlag, false, "h", "help")
	cmd.AddArg(&vFlag, false, "v", "version")
	cmd.AddArg(&dFlag, false, "d", "delete")
	cmd.AddArg(&pFlag, "", "p", "payload")
	cmd.AddArg(&lFlag, "", "l", "local")
	cmd.AddArg(&nFlag, "", "n", "notes")
	cmd.AddArg(&kFlag, "signingkeys.sigb", "k", "signingkeys")
	cmd.AddArg(&gFlag, "", "g", "gateway")
	cmd.AddArg(&sFlag, 3, "s", "shares")
	cmd.AddArg(&tFlag, 2, "t", "threshold")
	cmd.AddArg(&mFlag, 24, "m", "maxTtl")

	flag.Parse()

	if hFlag {
		flag.Usage()
		os.Exit(0)
	}

	if vFlag {
		fmt.Println(cmd.Version)
		os.Exit(0)
	}

	var hasPayloadOutput = false

	if pFlag == "" {
		fmt.Println("You need to specify a payload file")
		flag.Usage()
		os.Exit(1)
	}

	hasPayloadOutput = lFlag != "" || gFlag != ""

	if !hasPayloadOutput {
		fmt.Println("You need at least 1 payload output")
		flag.Usage()
		os.Exit(1)
	}

	plaPayload, err := fs.LoadFile(pFlag)
	cmd.ReportErrorAndExit(err)

	key, err := payload.NewArmoredPayloadKey("", nFlag)
	cmd.ReportErrorAndExit(err)

	encPayload, err := key.Encrypt(plaPayload)
	cmd.ReportErrorAndExit(err)

	if lFlag != "" {
		err = fs.SaveFile(encPayload, lFlag)
		cmd.ReportErrorAndExit(err)
	}

	if gFlag != ""{
		ipfs := remote.NewIpfsRemoteController(lFlag)
		rpi, err := ipfs.Add(encPayload)
		cmd.ReportErrorAndExit(err)
		key.IPFSAddress = rpi.Id
	}

	shares, err := key.Split(int(tFlag), int(sFlag))
	cmd.ReportErrorAndExit(err)

	signKeys, err := sign.NewSigningKeys()
	cmd.ReportErrorAndExit(err)

	signedShares := signKeys.SignShares(shares, uint64( mFlag))
	for i := 0; i < len(signedShares); i++ {
		fName := fmt.Sprintf("%s%s", signedShares[i].ShortId, cmd.ExtShare)
		jsonData, err := signedShares[i].Serialize()
		cmd.ReportErrorAndExit(err)
		fs.SaveFile(jsonData, fName)
	}

	if dFlag {
		err = fs.DeleteFile(pFlag)
		cmd.ReportErrorAndExit(err)
	}

	jsonSignKeys, err := signKeys.Serialize()
	cmd.ReportErrorAndExit(err)
	fs.SaveFile(jsonSignKeys, kFlag)
}
