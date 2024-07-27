package server

import (
	"chatroom/global"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		log.Println("模板解析错误", err)
		_, err = fmt.Fprint(w, "模板解析错误")
		if err != nil {
			log.Println("Fprint错误", err)
		}
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		log.Println("模板执行错误", err)
		if _, err = fmt.Fprint(w, "模板执行错误"); err != nil {
			log.Print("Fprint错误", err)
		}
	}
}
