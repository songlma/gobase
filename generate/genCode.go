package generate

import (
	"fmt"
)

/*
*
projectName 项目名称
appName 项目类型 web task daemon
service 具体的服务名称
*/
func CreateProject(projectFileDir, projectName, serviceType, serviceName string) {
	config := GenConfig{
		projectFileDir: projectFileDir,
		projectName:    projectName, //项目名称
		serviceType:    serviceType, //服务类型 例如 web task daemon
		serviceName:    serviceName, //服务名称  例如 sixhour
	}
	//通用
	config.GenFile("tmpl", "", false)
	config.GenFile("tmpl", "config", true)
	config.GenFile("tmpl", "ddl2struct", true)
	config.GenFile("tmpl/app", "constant", true)
	config.GenFile("tmpl/app", "dao", true)
	config.GenFile("tmpl/app", "helper", true)
	config.GenFile("tmpl/app", "logic", true)
	config.GenFile("tmpl/app", "model", true)
	switch serviceType {
	case "all":
		config.serviceType = "web"
		config.GenFile("tmpl", "dev", true)
		config.GenFile("tmpl/app/api", "web", true)
		config.GenFile("tmpl/app", "model", true)
		config.serviceType = "task"
		config.GenFile("tmpl", "dev", true)
		config.GenFile("tmpl/app/api", "task", true)
		config.serviceType = "daemon"
		config.GenFile("tmpl", "dev", true)
		config.GenFile("tmpl/app/api", "daemon", true)
	case "web":
		config.GenFile("tmpl", "dev", true)
		config.GenFile("tmpl/app/api", "web", true)
		config.GenFile("tmpl/app", "model", true)
	case "task":
		config.GenFile("tmpl", "dev", true)
		config.GenFile("tmpl/app/api", "task", true)
	case "daemon":
		config.GenFile("tmpl", "dev", true)
		config.GenFile("tmpl/app/api", "daemon", true)
	default:
		fmt.Println("appName need web or task or daemon")
	}
}
