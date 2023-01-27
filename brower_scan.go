package main

import "os"

// 扫描浏览器

const (
	DEFAULT_CHROME_USER_DATA_PATH = `\AppData\Local\Google\Chrome\User Data\`  // 默认的Chrome用户目录
	DEFAULT_EDGE_USER_DATA_PATH   = `\AppData\Local\Microsoft\Edge\User Data\` // 默认的Edge用户目录
)

// ScanBrowser 扫描浏览器
func ScanBrowser() []string {
	res := make([]string, 0)

	if chromePath := scanChromeBrowser(); len(chromePath) > 0 {
		res = append(res, chromePath)
	}

	return res
}

// scanChromeBrowser 扫描chrome浏览器
func scanChromeBrowser() string {
	path := os.Getenv("USERPROFILE") + DEFAULT_CHROME_USER_DATA_PATH
	if fp, fErr := os.Open(path); fErr == nil {
		fp.Close()
		return path
	}
	return ""
}
