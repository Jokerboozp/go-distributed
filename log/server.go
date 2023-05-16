package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

// 声明一个变量log指向*stlog.Logger类型的指针
var log *stlog.Logger

// 声明一个类型fileLog是一个字符串类型
type fileLog string

// 为fileLog类型添加Write方法
func (fl fileLog) Write(data []byte) (int, error) {

	//打开文件，如果不存在则创建，并且以追加的模式写入，权限为0600
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	//将data写入打开的文件中
	return f.Write(data)
}

// Run函数，用于初始化log
func Run(destination string) {
	log = stlog.New(fileLog(destination), "go ", stlog.LstdFlags)
}

// RegisterHandlers函数，用于注册http请求处理器
func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			//读取请求体
			msg, err := ioutil.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			//调用write函数写入日志
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// write函数，用于写入日志
func write(msg string) {
	log.Printf("%v\n", msg)
}
