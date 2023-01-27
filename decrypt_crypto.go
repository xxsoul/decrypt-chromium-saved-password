package main

import (
	"syscall"
	"unsafe"
)

// 调用win32 api解密，使用AES解密密文
const (
	CRYPTPROTECT_UI_FORBIDDEN = 0x1
)

var (
	dllCrypt32  = syscall.NewLazyDLL("crypt32.dll")
	dllKernel32 = syscall.NewLazyDLL("Kernel32.dll")

	procDecryptData = dllCrypt32.NewProc("CryptUnprotectData")
	procEncryptData = dllCrypt32.NewProc("CryptProtectData")
	procLocalFree   = dllKernel32.NewProc("LocalFree")
)

// Win32Encrypt 使用win32的加密方法
func Win32Encrypt(data []byte) ([]byte, error) {
	var res data_blob
	r, _, err := procEncryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, CRYPTPROTECT_UI_FORBIDDEN, uintptr(unsafe.Pointer(&res)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(res.pbData)))
	return res.toByteArr(), nil
}

// Win32Decrypt 使用win32的解密方法
func Win32Decrypt(data []byte) ([]byte, error) {
	var res data_blob
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, CRYPTPROTECT_UI_FORBIDDEN, uintptr(unsafe.Pointer(&res)))
	if r == 0 {
		return nil, err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(res.pbData)))
	return res.toByteArr(), nil
}

// data_blob win32 api要求的数据结构
type data_blob struct {
	cbData uint32
	pbData *byte
}

// 创建一个新的数据组
func newBlob(d []byte) *data_blob {
	if d == nil || len(d) < 1 {
		return &data_blob{}
	}

	return &data_blob{
		cbData: uint32(len(d)),
		pbData: &d[0],
	}
}

func (data *data_blob) toByteArr() []byte {
	res := make([]byte, data.cbData)
	copy(res, (*[1 << 15]byte)(unsafe.Pointer(data.pbData))[:])
	return res
}
