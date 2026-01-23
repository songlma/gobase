package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/songlma/gobase/generate"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("osGetwd():", err.Error())
		return
	}
	currentFileName := filepath.Base(path)

	args := os.Args
	help := "user gobase create projectName appName serviceName  eg:  gobase create usermsg daemon ql  gobase create appstory web"
	fmt.Println(fmt.Sprintf("args:%+v", args))
	if len(args) == 0 {
		fmt.Print("args size 0")
		return
	}
	if len(args) < 3 {
		fmt.Println(help)
		return
	}
	switch args[1] {
	case "create":
		if len(args) < 4 {
			fmt.Println(help)
			return
		}
		projectName := args[2]
		if matched, err := regexp.MatchString("^[a-z]*$", projectName); !matched || err != nil {
			fmt.Println("projectName must a-z")
			if err != nil {
				fmt.Println("err:", err.Error())
			}
			return
		}
		appName := args[3]
		if appName != "all" && appName != "web" && appName != "task" && appName != "daemon" {
			fmt.Println("appName need web or task or daemon")
		}

		serviceName := ""
		if len(args) >= 5 {
			serviceName = args[4]
		}
		if serviceName != "" {
			if matched, err := regexp.MatchString("^[a-z]*$", serviceName); !matched || err != nil {
				fmt.Println("serviceName must a-z")
				if err != nil {
					fmt.Println("err:", err.Error())
				}
				return
			}
		}
		//判断当前目录
		projectFileName := fmt.Sprintf("go-%s", projectName)
		projectFileDir := filepath.Join("./", projectFileName)
		if currentFileName == projectFileName {
			projectFileDir = filepath.Join("./")
		}
		generate.CreateProject(projectFileDir, projectName, appName, serviceName)
	case "help":
		fmt.Println(help)
	}
}
