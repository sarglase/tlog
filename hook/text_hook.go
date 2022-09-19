package hook

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type TextHook struct {
	// where to save your log
	filePath string
	// the number of files should be kept
	maxFileNum int
	writer     io.Writer
	// the format of the files name
	fileFormat string
}

var mux sync.Mutex
var currentDate = time.Now().Format("2006-01-02")
var defaultHook = TextHook{
	filePath:   "log",
	maxFileNum: 10,
	writer:     nil,
	fileFormat: "2006-01-02",
}

func New() *TextHook {
	return &defaultHook
}

func (th *TextHook) needNewFile() bool {
	now := time.Now().Format("2006-01-02")
	if now != currentDate {
		mux.Lock()
		defer mux.Unlock()
		currentDate = now
		return true
	}
	return false
}

func (th *TextHook) checkDir() error {
	if _, err := os.Stat(th.filePath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(th.filePath, os.ModePerm)
			if err != nil {
				fmt.Printf("mkdir log failed:%s\n", err.Error())
				return err
			}
			return nil
		} else {
			fmt.Printf("check dir exsist failed:%s\n", err.Error())
			return err
		}
	}
	return nil
}
func (th *TextHook) checkWriter() {
	if th.writer != nil {
		return
	}
	fileName := strings.TrimRight(th.filePath, "/") + "/" + fmt.Sprint(time.Now().Format(th.fileFormat)+".log")
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Sprintf("open file failed %s\n", err.Error())
		return
	} else {
		th.writer = f
	}
}

func (th *TextHook) SetConfig(filePath, fileFormat string, maxFileNum int) {
	th.fileFormat = fileFormat
	th.filePath = filePath
	th.maxFileNum = maxFileNum
}

func (th *TextHook) Write(message []byte) {
	if th.needNewFile() {
		fmt.Printf("neednew")
		th.writer = nil
	}
	err := th.checkDir()
	if err != nil {
		return
	}
	th.checkWriter()
	th.writer.Write(message)
}
