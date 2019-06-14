package cmdtool

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

/*
Office2pdf office文件转pdf，借助libreoffice
outdir 输出路径，绝对路径
srcfile 带转换文件路径，绝对路径
可借助docker使用，具体看dokcer file根据自己情况修改
*/
func Office2pdf(outdir, srcfile string) bool {
	cmds := "soffice --convert-to pdf --outdir " + outdir + " " + srcfile
	arr := strings.Split(cmds, " ")
	ar := arr[1:]
	//cmdtool := exec.Command("soffice","--convert-to","pdf","--outdir","存放路径"，"office文件路径" )
	cmd := exec.Command(arr[0], ar...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println(err)
		return false
	}
	fmt.Println("Result: " + out.String())
	return true
}
