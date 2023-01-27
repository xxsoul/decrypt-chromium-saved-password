package main

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DEFAULT_LOCAL_STATE = `Local State`
	DEFAULT_LOGIN_DATA  = `Default\Login Data`
)

var (
	logger = log.Default()
)

func main() {
	logger.Println("开始探测Chrome浏览器...")
	broPaths := ScanBrowser()
	if len(broPaths) < 1 {
		logger.Println("未探测到Chrome浏览器，退出")
		return
	}
	logger.Println("探测到Chrome浏览器，解密数据密钥...")
	encKey := loadEncKey(broPaths[0])
	encKeyHex := hex.EncodeToString(encKey)
	logger.Printf("解密数据密钥完毕，密钥：%s********%s\n", encKeyHex[:8], encKeyHex[len(encKeyHex)-8:])
	showSavedPass(broPaths[0]+DEFAULT_LOGIN_DATA, encKey, 5)
	logger.Println("Chrome浏览器数据处理完毕")
	logger.Print("按任意键继续...")
	// fmt.Scan()
	var ignore string
	fmt.Scan(&ignore)
}

// loadEncKey 加载加密密钥
func loadEncKey(path string) []byte {
	lsBytes, err := os.ReadFile(path + DEFAULT_LOCAL_STATE)
	if err != nil {
		return nil
	}

	localStateObj := struct {
		OsCrypt struct {
			EncKey string `json:"encrypted_key"`
		} `json:"os_crypt"`
	}{}
	if err = json.Unmarshal(lsBytes, &localStateObj); err != nil {
		logger.Fatalf("decode json data error:%v", err)
	}

	keyCipByte, err := base64.StdEncoding.DecodeString(localStateObj.OsCrypt.EncKey)
	if err != nil {
		logger.Fatalf("decode base64 data error:%v", err)
	}
	keyCipByte = keyCipByte[5:] // 去掉头部

	// 解密密钥
	keyPlaByte, err := Win32Decrypt(keyCipByte)
	if err != nil {
		logger.Fatalf("encryption key data decrypt error:%v", err)
	}
	return keyPlaByte
}

// showSavedPass 展示保存的密码
func showSavedPass(dbPath string, encKey []byte, count int) {
	logger.Println("加载用户数据库...")
	fTmp, err := os.OpenFile("tmp.db", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Fatalln("创建临时文件错误")
	}

	fDb, err := os.OpenFile(dbPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fTmp.Close()
		logger.Fatalln("读取用户数据库错误:" + err.Error())
	}

	if _, err = io.Copy(fTmp, fDb); err != nil {
		fDb.Close()
		fTmp.Close()
		logger.Fatalln("加载用户数据库错误:" + err.Error())
	}
	fDb.Close()
	fTmp.Close()

	logger.Println("开始读取用户数据库...")
	db, dbErr := sql.Open("sqlite3", "./tmp.db")
	if dbErr != nil {
		logger.Fatalln("读取用户数据库（2）错误:" + dbErr.Error())
	}

	rows, dbErr := db.Query(fmt.Sprintf("SELECT action_url, username_value, password_value FROM logins LIMIT %d", count*5))
	if dbErr != nil {
		logger.Fatalln("查询用户数据错误:" + dbErr.Error())
	}
	defer rows.Close()

	logger.Println("开始提取数据...")
	i := 0
	url, uname, psw := "", "", []byte{}
	for rows.Next() && i < count {
		url, uname, psw = "", "", []byte{}
		rows.Scan(&url, &uname, &psw)
		if (len(url) * len(uname) * len(psw)) == 0 {
			continue
		}
		i++

		// 解密密码
		iv := psw[3:15]
		cipData := psw[15:]
		plaData := aesGcmDecrypt(encKey, iv, cipData)
		pswStr := string(plaData)
		if len(pswStr) < 5 {
			logger.Printf("数据:%d\n网址:%s\n用户名:%s\n密码:**%s**\n", i, url, uname, pswStr)
		} else {
			logger.Printf("数据:%d\n网址:%s\n用户名:%s\n密码:%s****%s\n", i, url, uname, pswStr[:4], pswStr[len(pswStr)-4:])
		}
	}
	logger.Println("提取数据完毕，清理临时文件...")
	os.Remove("tmp.db")
	logger.Println("清理临时文件完毕")
}

// aesGcmDecrypt GCM解密数据
func aesGcmDecrypt(aesKey, aesIv, cipherData []byte) []byte {
	blocker, err := aes.NewCipher(aesKey)
	if err != nil {
		logger.Fatalln("加载数据密钥错误," + err.Error())
	}
	aead, err := cipher.NewGCM(blocker)
	if err != nil {
		logger.Fatalln("加载数据密钥参数错误," + err.Error())
	}

	res, err := aead.Open(nil, aesIv, cipherData, nil)
	if err != nil {
		logger.Println("解密数据错误," + err.Error())
		res = []byte("nil")
	}

	return res
}
