package mqListeners

import (
	"MQClient/logHelper"
	"MQClient/tools"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type CommonRequestQueue struct {
	*BaseListener
	logPrefix string
}

func(this *CommonRequestQueue) StartListen(done <-chan struct{}){
	this.logPrefix = "commonRequest"
	this.ListenByConsume(this.DoWork, done)
}

func(this *CommonRequestQueue) DoWork(msgs amqp.Delivery) error{
	prefix := this.logPrefix
	uid, err := uuid.NewV4()
	logHelper.FailOnError(prefix, err, "unknown uuid")
	logHelper.LogDefault(prefix, uid.String() + "|接收到数据[" + string(msgs.Body) + "]，开始排队...")
	var data = &CommonRequestData{}
	if err = json.Unmarshal(msgs.Body, &data); err != nil {
		logHelper.FailOnError(prefix, err, "failed to convert json data as [CommonRequestData]")
		msgs.Reject(true)
		return err
	}
	if response, err := tools.CommonRequest(data.Method, data.Url, data.Body, nil, "", "", data.Headers, ""); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|请求失败")
		msgs.Reject(true) //true 表示数据重新入列
	}else{
		logHelper.LogDefault(prefix, uid.String() + "|处理完成[" + response + "]")
		msgs.Ack(false)
	}
	return nil
}


//data mapping
type CommonRequestData struct {
	Url string	`json:"url"`
	Body string	`json:"body"`
	Headers map[string]string	`json:"headers"`
	Method string	`json:"method"`
}
