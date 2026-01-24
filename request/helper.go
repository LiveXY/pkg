package request

import (
	"bytes"
	"fmt"
	"sort"
)

type Header struct {
	Name  string
	Value string
}

type QueryParameter struct {
	Value any
	Key   string
}

func ConvertToQueryParams(params map[string]any) string {
	if params == nil {
		return ""
	}
	if len(params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for k, v := range params {
		if v == nil {
			continue
		}
		buffer.WriteString(fmt.Sprintf("%s=%v&", k, v))
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}
func ConvertToQueryParamsRepetition(params []QueryParameter) string {
	var buffer bytes.Buffer
	buffer.WriteString("?")
	for _, v := range params {
		if v.Value == nil {
			continue
		}
		buffer.WriteString(fmt.Sprintf("%s=%v&", v.Key, v.Value))
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}
func BuildTokenHeader(accessToken string) Header {
	return Header{
		Name:  "Authorization",
		Value: "Bearer " + accessToken,
	}
}

func Sort(data map[string]any) string {
	ksort := []string{}
	for k := range data {
		ksort = append(ksort, k)
	}
	sort.Strings(ksort)
	str := ""
	for _, k := range ksort {
		v := data[k]
		str += fmt.Sprintf("%s%v", k, v)
	}
	return str
}
