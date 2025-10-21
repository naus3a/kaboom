package fs

import(
	"os"
	"fmt"
	"crypto/rand"
)

// SaveFile saves data to a path
func SaveFile(data []byte, pth string)error{
	return os.WriteFile(pth, data, 0644)
}

// LoadFile loads a file from a path
func LoadFile(pth string)([]byte, error){
	return os.ReadFile(pth)
}

func DeleteFile(pth string)error{
	f, err := os.OpenFile(pth, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	fSz := info.Size()
	overwrite := make([]byte, fSz)
	_, err = rand.Read(overwrite)
	if err != nil {
		return err
	}
	_, err = f.Seek(0,0)
	if err != nil {
		return err
	}
	n, err := f.Write(overwrite)
	if err != nil {
                return err
        }
	if int64(n)!=fSz {
		return fmt.Errorf("wrong size during safe delete")
	}
	err = f.Sync()
	if err != nil {
                return err
        }
	return os.Remove(pth)
}
