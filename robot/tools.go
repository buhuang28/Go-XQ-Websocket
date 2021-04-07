package robot

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func DownPic(u string) []byte {
	u = strings.ReplaceAll(u,"\"","")
	request,_ := http.NewRequest("GET", u, nil)
	request.Header.Set("Upgrade-Insecure-Requests","1")
	request.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	//加入get参数
	q := request.URL.Query()
	request.URL.RawQuery = q.Encode()
	//isok := true
	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil || resp.Body == nil {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || body == nil {
		return nil
	}
	return body
}

func createLog() {
	logFileNmae := `./log/`+time.Now().Format("20060102")+".log"
	logFileAllPath := logFileNmae
	_,err :=os.Stat(logFileAllPath)
	exits := CheckFileIsExits(`log`)
	if !exits {
		_ = os.Mkdir("./log", os.ModePerm)
	}

	imgExist := CheckFileIsExits(`C:\images`)
	if !imgExist {
		_ = os.Mkdir(`C:\images`, os.ModePerm)
	}
	var f *os.File
	if  err != nil{
		f, _= os.Create(logFileAllPath)
	}else{
		//如果存在文件则 追加log
		f ,_= os.OpenFile(logFileAllPath,os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	logger = log.New(f, "", log.LstdFlags)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
}

func CheckFileIsExits(fileName string) bool {
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("GetOnlineQQs异常:",err)
			}else {
				writeFile("exc.txt","onStart异常")
			}
		}
	}()
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Contain(o string,c []string) bool{
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("GetOnlineQQs异常:",err)
			}else {
				writeFile("exc.txt","onStart异常")
			}
		}
	}()
	if c == nil || len(c) == 0 {
		return false
	}
	for _,v := range c {
		if o == v {
			return true
		}
	}
	return false
}

//追加
func writeFile(fileName,content string) bool {
	//if !checkFile(fileName) {
	//	return false
	//}
	defer func() {
		if err := recover(); err != nil { //产生了panic异常
			if logger != nil {
				logger.Println("GetOnlineQQs异常:",err)
			}else {
				writeFile("exc.txt","onStart异常")
			}
		}
	}()
	fd,_:=os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	buf:=[]byte(content)
	_, err := fd.Write(buf)
	fd.Close()
	if err == nil {
		return true
	}
	return false
}

var (
	imgReg = regexp.MustCompile(`\[Pic=(C.*?)\]`)
)

func CheckMsgImage(msg string) string {
	allString := imgReg.FindAllStringSubmatch(msg, -1)
	for _,v := range allString {
		if len(v) < 2 {
			logger.Println("检测到外星人的信息:",msg)
			continue
		}
		exits := CheckFileIsExits(v[1])
		if !exits {
			msg = strings.ReplaceAll(msg,v[1],"")
		}
	}
	return msg
}

func GetPwd() string {
	getwd, err := os.Getwd()
	logger.Println(getwd)
	if err != nil {
		logger.Println(err)
	}
	return getwd
}