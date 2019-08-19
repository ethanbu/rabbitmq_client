package mqListeners

import (
	"MQClient/logHelper"
	"MQClient/tools"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"strings"
)

type PayrollCalculationQueue struct {
	*BaseListener
	logPrefix string
}

func(this *PayrollCalculationQueue) StartListen(done <-chan struct{}) {
	this.logPrefix = "payroll_calc"
	this.ListenByConsume(this.DoWork, done)
}
func(this *PayrollCalculationQueue) DoWork(msgs amqp.Delivery) error {
	prefix := this.logPrefix
	uid, err := uuid.NewV4()
	logHelper.FailOnError(prefix, err, "unknown uuid")
	logHelper.LogDefault(prefix, uid.String() + "|接收到数据[" + string(msgs.Body) + "]，开始算薪...")
	data := &PayrollCalculationReciveData{}
	if err = json.Unmarshal(msgs.Body, &data); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|failed to convert json data as [PayrollCalculationReciveData]")
		msgs.Reject(true)
		return err
	}
	postData := generatePayrollCalculationData(data)
	header := map[string]string{"SOAPAction": `"http://tempuri.org/IStartPayrollCalculation/DoWork"`}
	postUrl := data.PAYROLL_SERVER
	postUrl = strings.TrimRight(postUrl, "?wsdl")
	if response, err := tools.PostData(postUrl, postData, nil, "", "", header, `text/xml; charset="UTF-8"`); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|请求失败")
		msgs.Reject(true)
	}else{
		logHelper.LogDefault(prefix, uid.String() + "|处理完成[" + response + "]")
		msgs.Ack(false)
	}
	return nil
}

func generatePayrollCalculationData(data *PayrollCalculationReciveData) string {
	params := make([]struct{
		Name string
		Value interface{}
	}, 7)
	addList := func(index int, k string, v interface{}){
		params[index] = struct {
			Name string
			Value interface{}
		}{
			Name: k,
			Value: v,
		}
	}
	addList(0,"companyID", data.CompanyId)
	addList(1,"taskID", data.TaskId)
	addList(2,"outerStepIndex", data.OuterStepIndex)
	addList(3,"instanceFlowID", data.InstanceFlowID)
	addList(4,"payrollGroupID", data.PayrollGroupID)
	addList(5,"calculationFlowID", data.CalculationFlowID)
	addList(6,"blnSubFlow", data.BlnSubFlow)
	return tools.GenerateSoapRequestData("DoWork", params)
}

type PayrollCalculationReciveData struct {
	TaskId string `json:"taskID"`
	CompanyId string `json:"companyID"`
	OuterStepIndex int	`json:"outerStepIndex"`
	InstanceFlowID int	`json:"instanceFlowID"`
	PayrollGroupID string	`json:"payrollGroupID"`
	CalculationFlowID string	`json:"calculationFlowID"`
	BlnSubFlow bool	`json:"blnSubFlow"`
	UserId string	`json:"userId"`
	CallBack string	`json:"callBack"`
	CallBack_P []string	`json:"callBack_P"`
	HTTP_REFERER string	`json:"HTTP_REFERER"`
	PAYROLL_SERVER string	`json:"PAYROLL_SERVER"`
}
