package generate

import "testing"

func TestGen(t *testing.T) {
	config := GenConfig{
		projectName:    "appnotice", //项目名称
		serviceType:    "task",      //服务类型 例如 web task daemon
		serviceName:    "sixhour",   //服务名称  例如 sixhour
		projectFileDir: "./tmp",
	}
	config.GenFile("tmpl/app/api", "daemon", true)
}

func TestGen1(t *testing.T) {
	Gen("appDaemon.txt", "/tmp/a/a.go", "./tmp", map[string]interface{}{"projectName": "go-test"}, true)
}
