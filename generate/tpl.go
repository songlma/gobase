package generate

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed tmpl/*
var templateFs embed.FS

type GenConfig struct {
	projectFileDir string //文件保存的目录
	projectName    string //项目名称 例如 apppush
	serviceType    string //服务类型 例如 web task daemon
	serviceName    string //服务名称  例如 sixhour
}

/*
*
parentPath 父级目录
currentDirPath 当前目录
projectFileDir 项目保存目录
dirNameFun
recursionDir 是否循环迭代目录下其他目录文件
*/
func (config GenConfig) GenFile(parentPath, currentDirPath string, recursionDir bool) {
	entries, _ := fs.ReadDir(templateFs, path.Join(parentPath, currentDirPath))
	for _, entry := range entries {
		if entry.IsDir() {
			if recursionDir {
				config.GenFile(path.Join(parentPath, currentDirPath), entry.Name(), recursionDir)
			}
		} else {
			fileName := entry.Name()
			//当前文件夹名称，dev特殊处理
			dirPath := currentDirPath
			if currentDirPath == "dev" {
				dirPath = config.serviceType
				if config.serviceName != "" {
					dirPath = config.serviceType + "-" + config.serviceName
				}
			}
			filepath := path.Join(config.projectFileDir, parentPath, dirPath, fileName)
			filepath = strings.ReplaceAll(filepath, ".tmpl", "")
			filepath = strings.ReplaceAll(filepath, "/tmpl", "")
			filepath = strings.ReplaceAll(filepath, "projectName", config.projectName)
			//模板参数
			var buildName = fmt.Sprintf("go-%s-%s", config.projectName, config.serviceType)
			if config.serviceName != "" {
				buildName = fmt.Sprintf("%s-%s", buildName, config.serviceName)
			}
			//例如appsearch-web 或者 reportorder-task-sixhour
			var project = fmt.Sprintf("%s-%s", config.projectName, config.serviceType)
			if config.serviceName != "" {
				project = fmt.Sprintf("%s-%s", project, config.serviceName)
			}
			//task类型入参
			taskName := ""
			if config.serviceType == "task" {
				taskName = " -task=test_task"
			}
			startApp := ""
			serviceName := ""
			switch config.serviceType {
			case "web":
				startApp = startWebAppTml
			case "task":
				startApp = startTaskAppTml
			case "daemon":
				startApp = startDaemonAppTml
			default:
				MustCheck(errors.New(fmt.Sprintf("appName %s main tml not found", config.serviceType)))
			}
			Gen(path.Join(parentPath, currentDirPath), entry.Name(), filepath, map[string]interface{}{
				"buildName":      buildName,
				"project":        project,
				"taskName":       taskName,
				"projectName":    config.projectName,
				"serviceType":    config.serviceType,
				"startApp":       startApp,
				"serviceName":    serviceName,
				"grpcServerName": strings.Title(config.projectName),
			}, strings.Contains(entry.Name(), ".go"))
		}
	}
}

const startWebAppTml = `
	addr := config.GetString("config.api.web_addr")
	webConf := web.NewAppConfig(addr)
	myapp = web.NewApp(ctx, webConf)
	if *taskName != "" {
		myapp.Once(ctx, *taskName)
		return
	}
	//启动主服务
	go myapp.Start(ctx)

	//监听信号
	listenSignal(ctx, myapp)

`

const startDaemonAppTml = `
	myapp = daemon.NewApp(serviceName)

	if *taskName != "" {
		myapp.Once(ctx, *taskName)
		return
	}
	addr := config.GetString("config.api.web_addr")
	defaultApp = app.NewDefaultApp(ctx, addr, "/inner", myapp)
	//启动监控服务
	go defaultApp.Start(ctx)
	//启动主服务
	go myapp.Start(ctx)
	//监听信号
	listenSignal(ctx, myapp)

`

const startTaskAppTml = `
	myapp = task.NewApp(ctx)
	myapp.Once(ctx, *taskName)

`

func Gen(templatePath, templateName, savePath string, data map[string]interface{}, goFile bool) {
	f, err := os.Open(savePath)
	if !(err != nil && os.IsNotExist(err)) {
		f.Close()
		return
	}
	generateTpl, err := template.ParseFS(templateFs, path.Join(templatePath, templateName))
	if err != nil {
		panic(err)
	}

	templates := generateTpl.Templates()
	for _, t := range templates {
		fmt.Println(t.ParseName)
	}
	var buf bytes.Buffer

	err = generateTpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		panic(err)
	}
	result := buf.Bytes()
	if goFile {
		result, err = format.Source(buf.Bytes())
		if err != nil {
			panic(err)
		}
	}
	dir := path.Dir(savePath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(savePath, result, 0644)
}
