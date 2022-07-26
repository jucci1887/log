/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 Date: 3/18/21 1:01 PM
*/
package services

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type SystemServices struct{}

// 获取当前执行程序的绝对目录路径
func GetCurrentDir() string {
	currentPath := CurrentAndAbsPath()
	return filepath.Dir(currentPath)
}

// 当前执行程序的绝对文件路径
func GetCurrentFilename() string {
	return CurrentAndAbsPath()
}

// 当前执行程序的绝对路径
func CurrentAndAbsPath() string {
	current := SetCurrentPath()
	return GetAbsPath(current)
}

// 设置当前执行程序的绝对路径
func SetCurrentPath() string {
	current := os.Args[0]
	path, err := exec.LookPath(current)
	if err != nil {
		log.Println("Set the current path error: ", err)
	}

	return path
}

// 获取当前执行程序的绝对路径
func GetAbsPath(current string) string {
	absPath, err := filepath.Abs(current)
	if err != nil {
		log.Println("Get the current absolute of path error: ", err)
	}

	return absPath
}

// 获取日志文件名
func GetLogsFilename() string {
	return "livestream.log"
}

// 获取日志文件路径
func GetLogsFilepath() string {
	logsPath := GetLogsDir()
	return filepath.Join(logsPath, GetLogsFilename())
}

// 获取日志文件内容前缀
func GetLogsPrefix() string {
	return ""
}

// 获取日志级别, 值为OFF则关闭日志
func GetLogsLevel() string {
	return "INFO"
}

// 获取进程文件路径
func GetPidPath() string {
	logsPath := GetLogsDir()
	return filepath.Join(logsPath, "livestream.pid")
}

// 获取客户端程序路径
func GetClientProgramPath() string {
	dir := GetCurrentDir()
	return filepath.Join(dir, string(os.PathSeparator), "client")
}

// 获取配置目录
func GetConfigDir() string {
	rootPath := GetRootPath()
	return filepath.Join(rootPath, "config", string(os.PathSeparator))
}

// 获取配置文件路径
func GetForwardConfigPath() string {
	configDir := GetConfigDir()
	return filepath.Join(configDir, "livestream.json")
}

// 获取配置
func GetForwardConfig() []byte {
	confPath := GetForwardConfigPath()
	config := FileService.GetFile(confPath)
	return config.ToByte()
}

// 获取日志目录
func GetLogsDir() string {
	rootPath := GetRootPath()
	return filepath.Join(rootPath, "logs", string(os.PathSeparator))
}

// 获取路径的上个目录
func GetLastPath(currentPath string) string {
	index := strings.LastIndex(currentPath, string(os.PathSeparator))
	return currentPath[:index]
}

// 获取项目根目录
func GetRootPath() string {
	dir := GetCurrentDir()
	rootPath := GetLastPath(dir)
	return filepath.Join(rootPath, string(os.PathSeparator))
}

// 环境变量解析: 根据环境变量的值替换字符串中的 ${var} or $var, 如果不存在任何环境变量, 则使用空字符串替换
func ParseEnvVar(varString string) string {
	return os.ExpandEnv(varString)
}

// 获取操作系统类型
func GetOS() string {
	return runtime.GOOS
}
