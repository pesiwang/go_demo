package main

import (
	"fmt"
	"strings"
	"text/template"
)

var temp string = `
###
input:
[
{{range .Msgs}}{{.}}{{end}}
]
output:
`

var msgs = []string{
	"\"msg1\",\n",
	"\"msg2\",\n",
	"\"msg3\"",
}

type InputData struct {
	Msgs []string
}

func main() {
	// 要注入的变量
	type Inventory struct {
		Material string
		Count    uint
	}
	sweaters := Inventory{"wool", 17}
	// 模板内容， {{.xxx}} 格式的都会被注入的变量替换
	text := `{{.Count}} items are made of {{.Material}}`

	result, err := ExecuteTemplate(text, sweaters)
	if err != nil {
		fmt.Printf("ExecuteTemplate failed: %v", err)
	}
	fmt.Printf("\n%v\n", result)

	inputData := InputData{Msgs: msgs}
	result, err = ExecuteTemplate(temp, inputData)
	if err != nil {
		fmt.Printf("ExecuteTemplate failed: %v", err)
	}
	fmt.Printf("\n%v\n", result)
}

func ExecuteTemplate(text string, data interface{}) (string, error) {
	// 初始化，解析
	tmpl, err := template.New("template").Parse(text)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil

}
