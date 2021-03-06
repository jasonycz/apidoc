// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build ignore

// 可通过在当前目录下手动执行 go run make.go
// 或是 go generate 来将模板内容打包到程序中。
//
// NOTE: 打包的代码中不能包含 ` 字符。

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
)

const (
	fileName    = "static.go" // 指定产生的文件名。
	packageName = "static"    // 指定包名。

	// 文件头部的警告内容
	warning = "// 该文件由 make.go 自动生成，请勿手动修改！\n\n"
)

// 指定所有需要序列化的文件名。
var assets = []string{
	"./style.css",
	"./app.js",
	"./index.html",
}

func main() {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := new(bytes.Buffer)
	w.WriteString(warning)

	// 输出包定义
	w.WriteString("package ")
	w.WriteString(packageName)
	w.WriteString("\n\n")

	makeStatic(w)

	data, err := format.Source(w.Bytes())
	if err != nil {
		panic(err)
	}

	l, err := f.Write(data)
	if err != nil {
		panic(err)
	}
	if l != len(data) {
		panic("未全部输出到文件")
	}
}

// 输出 assets 变量的整体。
func makeStatic(w *bytes.Buffer) {
	w.WriteString("var assets=map[string][]byte{\n")
	for _, file := range assets {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		w.WriteByte('"')
		w.WriteString(file)
		w.WriteString(`":`)
		w.WriteString("[]byte(`")
		w.Write(data)
		w.WriteString("`),")
	}
	w.WriteString("}\n")
}
