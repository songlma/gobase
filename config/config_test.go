package config

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var TestContext = context.Background()

func init() {

}

type Config struct {
	Dsn               string `json:"dsn"`
	Debug             bool   `json:"debug"`
	OpentracingPlugin bool   `json:"opentracing_plugin" mapstructure:"opentracing_plugin"`
}

func TestInit(t *testing.T) {
	//根据config path配置文件

	absPath, err := filepath.Abs(".")
	if err != nil {
		t.Error(err)
		return
	}

	Init(TestContext, []string{
		filepath.Join(absPath, "/config/config.yaml"),
		filepath.Join(absPath, "/secret/config.yaml"),
	}, true)
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			var conf Config
			err = UnmarshalKey("config.mysql.read", &conf)
			if err != nil {
				t.Error(err)
			}
			t.Log(fmt.Sprintf("%+v", conf))
			env := GetString("config.env")
			t.Log(fmt.Sprintf("env:%s", env))
			chuanglanNotice := GetInt("secret.chuanglan_notice")
			t.Log(fmt.Sprintf("chuanglanNotice:%d\n", chuanglanNotice))
		}
	}()
	group.Wait()
}
