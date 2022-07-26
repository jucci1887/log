/*
 Author: Kernel.Huang
 Mail: kernelman79@gmail.com
 Date: 3/18/21 1:01 PM
*/
package services

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type FileServices struct {
	name string
	Perm os.FileMode
}

var FileService = new(FileServices)

func (file *FileServices) FileInfo(filename string) os.FileInfo {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Println(err)
	}
	return fileInfo
}

func (file *FileServices) CheckFileActive(filename string) int64 {
	return file.FileInfo(filename).ModTime().Unix()
}

func (file *FileServices) GetFileSize(filename string) int64 {
	return file.FileInfo(filename).Size()
}

func (file *FileServices) GetFile(filename string) *FormatServices {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err.Error(), filename)
	}

	FormatService.to = buf
	return FormatService
}

func (file *FileServices) PutFile(filename string) *FileServices {
	file.name = filename
	return file
}

func (file *FileServices) Content(data interface{}) bool {
	switch value := data.(type) {

	case []byte:
		return file.FromByte(value)
	case string:
		return file.FromString(value)
	}

	return false
}

func (file *FileServices) FromByte(bytes []byte) bool {
	content := FormatService.FromByte(bytes).ToString()
	return file.BufIoPut(file.name, content)
}

func (file *FileServices) FromString(str string) bool {
	return file.BufIoPut(file.name, str)
}

func (file *FileServices) Check(err error) bool {
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func (file *FileServices) IoPut(filename string, content []byte) bool {
	err := ioutil.WriteFile(filename, content, file.Perm)
	return file.Check(err)
}

func (file *FileServices) BufIoPut(filename string, content string) bool {
	create := file.Create(filename)
	putFile := bufio.NewWriter(create)

	put, err := putFile.WriteString(content)
	if err != nil {
		log.Printf("Write data error")
		return false
	}

	err = putFile.Flush()
	if err != nil {
		log.Printf("Push data error")
		return false
	}

	defer create.Close()
	log.Printf("Write %d bytes\n", put)

	return true
}

func (file *FileServices) OpenAndAppend(filename string) *os.File {
	fOpen, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, file.Perm)
	file.Check(err)
	defer fOpen.Close()
	return fOpen
}

func (file *FileServices) OpenAndOverwrite(filename string) (*os.File, error) {
	open, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, file.Perm)
	if err != nil {
		return nil, err
	}

	return open, nil
}

func (file *FileServices) Create(filename string) *os.File {
	fCreate, err := os.Create(filename)
	file.Check(err)
	return fCreate
}

/**
 * 检查文件或目录是否存在
 */
func (file *FileServices) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 创建目录
 */
func (file *FileServices) CreateDir(path string, perm os.FileMode) bool {
	err := os.Mkdir(path, perm)
	if err != nil {
		return file.Check(err)
	}
	return true
}

/**
 * 备份文件
 * path: source file
 * backupDir e.g: "/tmp/recycle/"
 */
func (file *FileServices) BackupFile(path string, backupDir string) bool {
	exist, err := file.PathExists(backupDir)
	if err != nil {
		file.Check(err)
	}

	if !exist {
		mk := file.CreateDir(backupDir, 0744)
		if mk {
			return file.MoveFile(path, backupDir)
		}
		return false
	}
	return file.MoveFile(path, backupDir)
}

/**
 * 移动文件
 */
func (file *FileServices) MoveFile(oldPath string, newPath string) bool {
	filename := StringService.SplitFileName(oldPath)
	exist, err := file.PathExists(newPath + filename)
	if err != nil {
		file.Check(err)
	}

	if exist {
		err = os.Rename(newPath+filename, newPath+filename+".bak"+FormatService.FromInt64(time.Now().Unix()).ToString())
		if err != nil {
			file.Check(err)
		}
	}

	err = os.Rename(oldPath, newPath+filename)
	if err != nil {
		return file.Check(err)
	}

	err = os.Chmod(newPath+filename, 0644)
	if err != nil {
		file.Check(err)
	}
	return true
}
