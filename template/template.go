package template

import (
	"bytes"
	htmltemplate "html/template"
	"path"
	"strings"
	texttemplate "text/template"

	"github.com/livexy/pkg/crypto/md5x"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/valyala/fasttemplate"
)

var fastTemplateMap = cmap.New[*fasttemplate.Template]()
var textTemplateMap = cmap.New[*texttemplate.Template]()
var htmlTemplateMap = cmap.New[*htmltemplate.Template]()
var textkeys = []string { "{{if ", "{{ if ", "{{else}}", "{{ else }}", "{{end}}", "{{ end }}" }
var htmlkeys = []string { "</html>", "</body>" }

func check(text string, keys []string) bool {
	var flag bool
	for _, v := range textkeys {
		if strings.Contains(text, v) {
			flag = true
			break
		}
	}
	return flag
}

func FastTemplate(text string, param map[string]any) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	if check(text, htmlkeys) {
		return HtmlTemplate(text, param)
	}
	if check(text, textkeys) {
		return TextTemplate(text, param)
	}
	key := md5x.MD5(text)
	obj, ok := fastTemplateMap.Get(key)
	if ok {
		return obj.ExecuteString(param), nil
	} else {
		tmpl, err := fasttemplate.NewTemplate(text, "{{.", "}}")
		if err != nil {
			return "", err
		}
		fastTemplateMap.Set(key, tmpl)
		return tmpl.ExecuteString(param), nil
	}
}

func TextTemplate(text string, param map[string]any) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	key := md5x.MD5(text)
	obj, ok := textTemplateMap.Get(key)
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
		textTemplateMap.Set(key, tmpl)
		return buf.String(), nil
	}
}

func HtmlTemplate(text string, param map[string]any) (string, error) {
	if len(text) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	key := md5x.MD5(text)
	obj, ok := htmlTemplateMap.Get(key)
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
		htmlTemplateMap.Set(key, tmpl)
		return buf.String(), nil
	}
}

func HtmlFileTemplate(filename string, param interface{}, funcMap map[string]interface{}) (string, error) {
	name := path.Base(filename)
	var buf bytes.Buffer
	obj, ok := htmlTemplateMap.Get(filename)
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
		htmlTemplateMap.Set(filename, tmpl)
		return buf.String(), nil
	}
}
