package template

import (
	"bytes"
	htmltemplate "html/template"
	"path"
	"strings"
	texttemplate "text/template"

	"github.com/livexy/pkg/crypto/md5x"

	"github.com/puzpuzpuz/xsync/v4"
	"github.com/valyala/fasttemplate"
)

var fastTemplateMap = xsync.NewMap[string, *fasttemplate.Template]()
var textTemplateMap = xsync.NewMap[string, *texttemplate.Template]()
var htmlTemplateMap = xsync.NewMap[string, *htmltemplate.Template]()
var textkeys = []string{"{{if ", "{{ if ", "{{else}}", "{{ else }}", "{{end}}", "{{ end }}"}
var htmlkeys = []string{"</html>", "</body>"}

func check(text string, keys []string) bool {
	var flag bool
	for _, v := range keys {
		if strings.Contains(text, v) {
			flag = true
			break
		}
	}
	return flag
}

// FastTemplate 使用 fasttemplate 高效处理简单占位符替换，自动适配 HTML 和 text 模式
func FastTemplate(text string, param map[string]any) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	if check(text, htmlkeys) {
		return HTMLTemplate(text, param)
	}
	if check(text, textkeys) {
		return TextTemplate(text, param)
	}
	key := md5x.MD5(text)
	obj, ok := fastTemplateMap.Load(key)
	if ok {
		return obj.ExecuteString(param), nil
	} else {
		tmpl, err := fasttemplate.NewTemplate(text, "{{.", "}}")
		if err != nil {
			return "", err
		}
		fastTemplateMap.Store(key, tmpl)
		return tmpl.ExecuteString(param), nil
	}
}

// TextTemplate 使用标准库 text/template 处理模板
func TextTemplate(text string, param map[string]any) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	key := md5x.MD5(text)
	obj, ok := textTemplateMap.Load(key)
	if ok {
		err := obj.Execute(&buf, param)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	} else {
		tmpl, err := texttemplate.New(key).Parse(text)
		if err != nil {
			return "", err
		}
		err = tmpl.Execute(&buf, param)
		if err != nil {
			return "", err
		}
		textTemplateMap.Store(key, tmpl)
		return buf.String(), nil
	}
}

// HTMLTemplate 使用标准库 html/template 处理 HTML 模板（自动转义安全）
func HTMLTemplate(text string, param map[string]any) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	key := md5x.MD5(text)
	obj, ok := htmlTemplateMap.Load(key)
	if ok {
		err := obj.Execute(&buf, param)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	} else {
		tmpl, err := htmltemplate.New(key).Parse(text)
		if err != nil {
			return "", err
		}
		err = tmpl.Execute(&buf, param)
		if err != nil {
			return "", err
		}
		htmlTemplateMap.Store(key, tmpl)
		return buf.String(), nil
	}
}

// HTMLFileTemplate 从文件加载并处理 HTML 模板
func HTMLFileTemplate(filename string, param any, funcMap map[string]any) (string, error) {
	name := path.Base(filename)
	var buf bytes.Buffer
	obj, ok := htmlTemplateMap.Load(filename)
	if ok {
		err := obj.ExecuteTemplate(&buf, name, param)
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	} else {
		tmpl, err := htmltemplate.New(name).Funcs(funcMap).ParseFiles(filename)
		if err != nil {
			return "", err
		}
		err = tmpl.ExecuteTemplate(&buf, name, param)
		if err != nil {
			return "", err
		}
		htmlTemplateMap.Store(filename, tmpl)
		return buf.String(), nil
	}
}
