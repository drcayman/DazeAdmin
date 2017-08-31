package main

import (
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/crabkun/DazeAdmin/database"
	"github.com/crabkun/DazeAdmin/console"
)

type S_config struct{
	DatabaseDriver string
	DatabaseConnectionString string
}
var config S_config
func main(){
	fmt.Println("DazeAdmin v1.0 DazeProxyV3数据库管理工具")
	buf,ReadErr:=ioutil.ReadFile("config.json")
	if ReadErr!=nil{
		fmt.Println("配置文件(config.json)读取错误："+ReadErr.Error())
		os.Exit(-1)
	}
	JsonErr:=json.Unmarshal(buf,&config)
	if JsonErr!=nil{
		fmt.Println("配置文件格式错误！请参考DefaultConfig.json并严格按照JSON格式来修改config.json(",JsonErr.Error(),")")
		os.Exit(-1)
	}
	fmt.Println("配置文件读取成功")
	os.Stdout.Sync()
	database.LoadDatabase(config.DatabaseDriver,config.DatabaseConnectionString)
	console.Start()
}