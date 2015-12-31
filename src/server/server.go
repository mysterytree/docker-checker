package main

import (
	"cstructs"
	"encoding/json"
	"fmt"
	"github.com/franela/goreq"
	"io/ioutil"
	"net/http"
)

var (
	token       string
	channelName string
	userName    string
	port        string
)

func main() {
	token = cstructs.GetEnv("SERVER_Token", "https://hooks.slack.com/services/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")
	channelName = cstructs.GetEnv("SERVER_CHANNEL", "#test")
	userName = cstructs.GetEnv("SERVER_USERNAME", "BombChecker")
	port = cstructs.GetEnv("SERVER_PORT", "50075")

	//设置访问路由
	http.HandleFunc("/", index)
	//设置访问端口
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	//解析参数，默认是不会解析的
	r.ParseForm()
	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		//解析
		var container cstructs.Container
		json.Unmarshal([]byte(result), &container)
		//处理
		statusFilter(container)
	} else {
		fmt.Printf("错误请求")
	}
}

func statusFilter(container cstructs.Container) {
	//颜色
	var color string
	//区分颜色
	switch container.Status {
	case "start":
		color = "good"
	case "die":
		color = "danger"
	}
	//组装message
	if color != "" {
		//组装attachment
		var attachments []cstructs.Attachment
		//组装field
		var fields []cstructs.Field
		alertLog := cstructs.Field{"Log", container.Log, false}
		alertTime := cstructs.Field{"Time", container.Time, true}
		fields = append(fields, alertLog)
		fields = append(fields, alertTime)
		//attachment
		attachment := cstructs.Attachment{
			Color:   color,
			Text:    " At " + container.Time,
			Pretext: container.Name + " " + container.Status + " !",
			Fields:  fields,
		}
		//将attachment放入组
		attachments = append(attachments, attachment)
		//最终的message
		message := cstructs.Message{
			Channel:     channelName,
			Username:    userName,
			Attachments: attachments,
		}
		//发送
		sendSlack(&message)
	}

}

func sendSlack(message *cstructs.Message) {
	//发送
	res, err := goreq.Request{
		Method: "POST",
		Uri:    token,
		Body:   message,
	}.Do()
	if err != nil {
		fmt.Println(err)
	}
	//结果
	resString, err := res.Body.ToString()
	fmt.Printf("result： %s \n", resString)
}
