package fs

import(
	"bytes"
	"testing"
)

func TestFileOps(t *testing.T){
	const testFilePth = "test.txt"
	testData := []byte("this is a test")
	_ = DeleteFile(testFilePth)
	var err error = nil 
	defer func(){
		err = DeleteFile(testFilePth)
		if err != nil {
			t.Errorf("FAILED: could not delete: %v", err)
		}else{
			t.Log("file deleted")
		}
	}()
	err = SaveFile(testData, testFilePth)
	if err != nil {
		t.Errorf("FAILED: could not save: %v", err)
	}else{
		t.Log("file saved")
	}
	loadedData, err := LoadFile(testFilePth)
	if err!= nil {
		t.Errorf("FAILED: could not load: %v", err)
	}else{
		t.Log("file loaded")
	}
	if !bytes.Equal(testData, loadedData){
		t.Errorf("FAILED: expected %s, got %s", testData, loadedData)
	}else{
		t.Log("contenr succesfully checked")
	}
}
