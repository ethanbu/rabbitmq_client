package mqListeners

import (
	"MQClient/logHelper"
	"MQClient/tools"
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"strings"
)

type ImportExcelListener struct {
	*BaseListener
	logPrefix string
}

func(this *ImportExcelListener) StartListen(done <-chan struct{}){
	this.logPrefix = "importExcel"
	this.ListenByConsume(this.DoWork, done)
}

func(this *ImportExcelListener) DoWork(msgs amqp.Delivery) error {
	var requestData = ImportExcelUploadRequestData{}
	prefix := this.logPrefix
	uid, err := uuid.NewV4()
	logHelper.FailOnError(prefix, err, "unknown uuid")
	logHelper.LogDefault(prefix, uid.String() + "|接收到数据" + string(msgs.Body) + "]，开始排队进行导入...")
	if err = json.Unmarshal(msgs.Body, &requestData); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|failed to convert json data as [ImportExcelUploadRequestData]")
		msgs.Reject(true) //true 表示数据重新入列
		return err
	}
	postUrl := "/rest.php/c/common/Upload"
	requestType := "http://"
	referer := requestData.HTTPREFERER
	if strings.HasPrefix(referer, "https://") {
		requestType = "https://"
	}
	referer = strings.TrimLeft(referer, requestType)
	referer = referer[0: strings.Index(referer, "/")]
	postUrl = requestType + referer + postUrl
	if response, err := tools.PostData(postUrl, string(msgs.Body), nil, "", "", nil, ""); err != nil {
		logHelper.FailOnError(prefix, err, uid.String() + "|请求失败")
		msgs.Reject(true) //true 表示数据重新入列
	}else{
		logHelper.LogDefault(prefix, uid.String() + "|处理完成[" + response + "]")
		msgs.Ack(false)
	}
	return nil
}

type ImportExcelUploadRequestData struct {
	Data ImportExcelUploadedInfo	`json:"data"`
	UserId string `json:"user_id"`
	CompanyId string `json:"company_id"`
	Method string	`json:"method"`
	CronjobInsId string	`json:"cronjobInsId"`
	HTTPREFERER string	`json:"HTTP_REFERER"`
	ClientCode string	`json:"client_code"`
}

type ImportExcelUploadedInfo struct {
	Method string	`json:"method"`
	Data []ImportExcelFieldsMap	`json:"data"`
	Location string	`json:"location"`
	TableName string	`json:"tableName"`
	UploadType string	`json:"uploadType"`
	WebFile string	`json:"webFile"`
	HeaderRow string	`json:"headerRow"`
	FirstRow string	`json:"firstRow"`
	LastRow string	`json:"lastRow"`
	Options map[string] string	`json:"options"`
	DataService string	`json:"dataService"`
	SheetName string	`json:"sheetName"`
}
type ImportExcelFieldsMap struct {
	SourceField string	`json:"sourceField"`
	TargetField string	`json:"targetField"`
	TargetName string	`json:"targetName"`
	RefFilterTable string	`json:"refFilterTable"`
	RefFilterField string	`json:"refFilterField"`
	RefTable string	`json:"refTable"`
	RefField string	`json:"refField"`
	TableCode string	`json:"tableCode"`
	FieldType string	`json:"fieldType"`
	MultiSelect string	`json:"multiSelect"`
	Is_required string	`json:"is_required"`
}