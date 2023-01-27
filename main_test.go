package main

import (
	"encoding/base64"
	"testing"
)

func TestEncDec(t *testing.T) {
	t.Log("begin")
	data := []byte("hello world")

	encData, err := Win32Encrypt(data)
	if err != nil {
		t.Fatalf("encrypt error:%v", err)
	}
	t.Logf("enc data is:%s", base64.StdEncoding.EncodeToString(encData))
	decData, err := Win32Decrypt(encData)
	if err != nil {
		t.Fatalf("decrypt error:%v", err)
	}
	t.Logf("dec data is:%s", string(decData))

	t.Log("done!")
}

func TestScanBrowser(t *testing.T) {
	t.Log("begin")

	path := ScanBrowser()
	for _, s := range path {
		t.Log(s)
	}

	t.Log("done!")
}
