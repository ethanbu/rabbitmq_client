package tools

import (
	"fmt"
	"reflect"
	"strconv"
)

func GenerateSoapRequestData(actionName string, params []struct{
	Name string
	Value interface{}
}) string {
	retStr := `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">`
	retStr += fmt.Sprintf(`<Body><%s xmlns="http://tempuri.org/">`, actionName)
	for _, item := range params {
		var strV string
		rv := reflect.ValueOf(item.Value)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := rv.Int()
			strV = strconv.FormatInt(i, 10)
		case reflect.Bool:
			strV = strconv.FormatBool(rv.Bool())
		default:
			strV = item.Value.(string)
		}
		retStr += fmt.Sprintf("<%s>%s</%s>", item.Name, strV, item.Name)
	}
	retStr += fmt.Sprintf(`</%s></Body></Envelope>`, actionName)

	return retStr
}
