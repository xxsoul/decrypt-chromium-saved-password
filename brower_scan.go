package main

import "os"

// 扫描浏览器

const (
	DEFAULT_CHROME_USER_DATA_PATH = `\AppData\Local\Google\Chrome\User Data\`  // 默认的Chrome用户目录
	DEFAULT_EDGE_USER_DATA_PATH   = `\AppData\Local\Microsoft\Edge\User Data\` // 默认的Edge用户目录
)

// browserInfo 浏览器信息
type browserInfo struct {
	name     string // 浏览器名称
	userPath string // 浏览器路径
}

// ScanBrowser 扫描浏览器
func ScanBrowser() []browserInfo {
	res := make([]browserInfo, 0)

	if browInf := scanChromeBrowser(); browInf != nil {
		res = append(res, *browInf)
	}

	if browInf := scanEdgeBrowser(); browInf != nil {
		res = append(res, *browInf)
	}

	return res
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
