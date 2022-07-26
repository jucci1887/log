/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 Date: 3/18/21 1:01 PM
*/
package log

import (
	"fmt"
	"github.com/kavanahuang/log/services"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const DateFormat = "2006-01-02"

type LEVEL byte

const (
	TRACE LEVEL = iota
	INFO
	WARN
	ERROR
	OFF
)

var fileLog *LogServices

type LoggerConf struct {
	FileDir  string
	FileName string
	Prefix   string
	Level    string
}

type LogServices struct {
	fileDir  string
	fileName string
	prefix   string
	date     *time.Time
	logFile  *os.File
	lg       *log.Logger
	logLevel LEVEL
	mu       *sync.RWMutex
	logChan  chan string
}

// 初始化日志配置
func BootLogger() (err error) {
	conf := &LoggerConf{
		FileDir:  services.GetLogsDir(),
		FileName: services.GetLogsFilename(),
		Prefix:   services.GetLogsPrefix(),
		Level:    services.GetLogsLevel(),
	}

	f := &LogServices{
		fileDir:  conf.FileDir,
		fileName: conf.FileName,
		prefix:   conf.Prefix,
		mu:       new(sync.RWMutex),
		logChan:  make(chan string, 8000),
	}

	if strings.EqualFold(conf.Level, "OFF") {
		f.logLevel = OFF
	} else if strings.EqualFold(conf.Level, "TRACE") {
		f.logLevel = TRACE
	} else if strings.EqualFold(conf.Level, "WARN") {
		f.logLevel = WARN
	} else if strings.EqualFold(conf.Level, "ERROR") {
		f.logLevel = ERROR
	} else {
		f.logLevel = INFO
	}

	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	f.date = &t

	if f.isMustSplit() {
		if err = f.split(); err != nil {
			return
		}
	} else {
		f.isExistOrCreate()

		logFile := filepath.Join(f.fileDir, f.fileName)

		f.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return
		}

		f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)
	}

	go f.logWriter()
	go f.fileMonitor()

	fileLog = f
	return
}

// 日志文件是否分割
func (f *LogServices) isMustSplit() bool {
	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	return t.After(*f.date)
}

// 检查日志文件是否存在，不存在则创建
func (f *LogServices) isExistOrCreate() {
	_, err := os.Stat(f.fileDir)
	if err != nil && !os.IsExist(err) {
		mkdirErr := os.Mkdir(f.fileDir, 0755)
		if mkdirErr != nil {
			log.Println("Create dir failed, error: ", mkdirErr)
		}
	}
}

// 分割日志
func (f *LogServices) split() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	logFile := filepath.Join(f.fileDir, f.fileName)
	logFileBak := logFile + "." + f.date.Format(DateFormat)

	if f.logFile != nil {
		_ = f.logFile.Close()
	}

	err = os.Rename(logFile, logFileBak)
	if err != nil {
		return
	}

	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	f.date = &t

	f.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)
	return
}

// 日志写入
func (f *LogServices) logWriter() {
	defer func() { recover() }()

	for {
		str := <-f.logChan

		f.mu.RLock()
		_ = f.lg.Output(2, str)
		f.mu.RUnlock()
	}
}

// 日志分割监控
func (f *LogServices) fileMonitor() {
	defer func() { recover() }()

	timer := time.NewTicker(30 * time.Second)
	for {
		<-timer.C

		if f.isMustSplit() {
			if err := f.split(); err != nil {
				f.Error("Log split error: %v\n", err)
			}
		}
	}
}

// 关闭日志
func CloseLogger() {
	if fileLog != nil {
		close(fileLog.logChan)
		fileLog.lg = nil
		_ = fileLog.logFile.Close()
	}
}

// 输出格式化日志
func Printf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fileLog.logChan <- fmt.Sprintf("[%v:%v]", fmt.Sprintf(format, v...)+filepath.Base(file), line)
}

// 输出格式化日志
func Print(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fileLog.logChan <- fmt.Sprintf("[%v:%v]", fmt.Sprint(v...)+filepath.Base(file), line)
}

// 输出格式化日志
func Println(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fileLog.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintln(v...)
}

func Fatally(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	fileLog.logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[ERROR] [")+filepath.Base(file), line) + fmt.Sprintln(v...)
	_ = log.Output(2, fmt.Sprintln(v))
	os.Exit(1)
}

// 输出跟踪日志
func Trace(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	if fileLog.logLevel <= TRACE {
		fileLog.logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[TRACE] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出信息日志
func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	s := fmt.Sprintf("%v:%v:%v%v]", fmt.Sprintf("[INFO] [")+filepath.Base(file), line, format, v)
	fmt.Printf("\033[0;40;32m%s\033[0m\n", s)
	if fileLog.logLevel <= INFO {
		fileLog.logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[INFO] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出警告日志
func Warn(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLog.logLevel <= WARN {
		fileLog.logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[WARN] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}

// 输出错误日志
func Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLog.logLevel <= ERROR {
		fileLog.logChan <- fmt.Sprintf("%v:%v]", fmt.Sprintf("[ERROR] [")+filepath.Base(file), line) + fmt.Sprintf(" "+format, v...)
	}
}
