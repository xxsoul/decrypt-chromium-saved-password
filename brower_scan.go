package main

import "os"

// 扫描浏览器

const (
	DEFAULT_CHROME_USER_DATA_PATH = `\AppData\Local\Google\Chrome\User Data\`  // 默认的Chrome用户目录
	DEFAULT_EDGE_USER_DATA_PATH   = `\AppData\Local\Microsoft\Edge\User Data\` // 默认的Edge用户目录
)

var (
	chromiumBrowPathMap = map[string]string{
		"Chrome":    `\AppData\Local\Google\Chrome\User Data\`,
		"Edge":      `\AppData\Local\Microsoft\Edge\User Data\`,
		"Chromium":  `\AppData\Local\Chromium\User Data\`,
		"Opera":     `\AppData\Roaming\Opera Software\Opera Stable\`,
		"Vivaldi":   `\AppData\Local\Vivaldi\User Data\`,
		"Coccoc":    `\AppData\Local\CocCoc\Browser\User Data\`,
		"Brave":     `\AppData\Local\BraveSoftware\Brave-Browser\User Data\`,
		"Yandex":    `\AppData\Local\Yandex\YandexBrowser\User Data\`,
		"360Speed":  `\AppData\Local\360chrome\Chrome\User Data\`,
		"QQ":        `\AppData\Local\Tencent\QQBrowser\User Data\`,
		"Dcbrowser": `\AppData\Local\DCBrowser\User Data\User Data\`,
		"Sougou":    `\AppData\Roaming\SogouExplorer\Webkit\`,
	}

	firefoxBrowPathMap = map[string]string{
		"Firefox": `\AppData\Roaming\Mozilla\Firefox\Profiles`,
	}
)

// browserInfo 浏览器信息
type browserInfo struct {
	name     string // 浏览器名称
	userPath string // 浏览器路径
}

// ScanBrowser 扫描浏览器
// func ScanBrowser() []browserInfo {
// 	res := make([]browserInfo, 0)

// 	if browInf := scanChromeBrowser(); browInf != nil {
// 		res = append(res, *browInf)
// 	}

// 	if browInf := scanEdgeBrowser(); browInf != nil {
// 		res = append(res, *browInf)
// 	}

// 	return res
// }

func ScanChromiumBrowser() []browserInfo {
	browList := make([]browserInfo, 0)
	homePath := os.Getenv("USERPROFILE")
	for name, broPath := range chromiumBrowPathMap {
		path := homePath + broPath
		if fp, fErr := os.Open(path); fErr == nil {
			fp.Close()
			browList = append(browList, browserInfo{
				name:     name,
				userPath: path,
			})
		}
	}
	return browList
}

func ScanFirefoxBrowser() []browserInfo {
	browList := make([]browserInfo, 0)
	homePath := os.Getenv("USERPROFILE")
	for name, broPath := range firefoxBrowPathMap {
		path := homePath + broPath
		if fp, fErr := os.Open(path); fErr == nil {
			fp.Close()
			browList = append(browList, browserInfo{
				name:     name,
				userPath: path,
			})
		}
	}
	return browList
}

// scanChromeBrowser 扫描chrome浏览器
func scanChromeBrowser() *browserInfo {
	path := os.Getenv("USERPROFILE") + DEFAULT_CHROME_USER_DATA_PATH
	if fp, fErr := os.Open(path); fErr == nil {
		fp.Close()
		return &browserInfo{
			name:     "Chrome",
			userPath: path,
		}
	}
	return nil
}

// scanEdgeBrowser 扫描edge浏览器
func scanEdgeBrowser() *browserInfo {
	path := os.Getenv("USERPROFILE") + DEFAULT_EDGE_USER_DATA_PATH
	if fp, fErr := os.Open(path); fErr == nil {
		fp.Close()
		return &browserInfo{
			name:     "Edge",
			userPath: path,
		}
	}
	return nil
}
