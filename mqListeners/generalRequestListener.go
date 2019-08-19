package mqListeners

import (
	"MQClient/logHelper"
	"MQClient/tools"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"strings"
	"time"
)

type GeneralRequestQueue struct {
	*BaseListener
	logPrefix string
}

func(this *GeneralRequestQueue) StartListen(done <-chan struct{}) {
	this.logPrefix = "general_request"
	this.ListenByConsume(this.DoWork, done)
}

func(this *GeneralRequestQueue) DoWork(msgs amqp.Delivery) error {
	prefix := this.logPrefix
	uid, err := uuid.NewV4()
	logHelper.FailOnError(prefix, err, "unknown uuid")
	logHelper.LogDefault(prefix, uid.String() + "|接收到数据[" + string(msgs.Body) + "]，开始排队...")
	data := &GeneralData{}
	if err = json.Unmarshal(msgs.Body, &data); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|failed to convert json data as [GeneralData]")
		msgs.Reject(true)
		return err
	}
	referer := strings.ToLower(data.Referer)
	requestType := "http://"
	if strings.HasPrefix(referer, "https://") {
		requestType = "https://"
	}
	host := strings.Replace(referer, requestType, "", 1)
	host = host[0: strings.Index(host, "/")]
	postUrl := requestType + host + data.CallBack
	if strings.Index(strings.ToLower(data.CallBack), "sendemailnow_mq") > -1 {
		time.Sleep(time.Second * 3)
	}
	jsonArgsV, err := json.Marshal(data.Msg)
	if err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|failed to convert json data as [jsonArgsV]")
		return err
	}
	json_args := map[string]interface{}{"json_args": string(jsonArgsV)}
	argsJsonV, err := json.Marshal(json_args)
	if err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|failed to convert json data as [argsJsonV]")
		return err
	}
	postData := "argsJson=" + string(argsJsonV)
	if response, err := tools.PostData(postUrl, postData, nil, host, data.Referer, nil, ""); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|请求失败")
		msgs.Reject(true)
	}else{
		logHelper.LogDefault(prefix, uid.String() + "|处理完成[" + response + "]")
		msgs.Ack(false)
	}
	return nil
}

type GeneralData struct {
	Msg interface{}	`json:"msg"`
	CallBack string 	`json:"callback"`
	Referer string 	`json:"referer"`
}
