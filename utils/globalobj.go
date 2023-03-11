package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

/*
存储一切有关gzinx框架的全局擦书，供其他模块使用
一些参数是可以通过gzinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前gzinx全局的server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		GZinx
	*/
	Version        string //当前gzinx的版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前gzinx框架数据包的最大值
}

/*
定义一个全局的对外Globalobj
*/

var GlobalObject *GlobalObj

/*
从gzinx.json去加载用户自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/gzinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到struct中
	json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
提供一个init方法，初始化当前的GlobalObject
*/
func init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "GZinxServerApp",
		Host:           "0.0.0.0",
		Version:        "V0.5",
		TcpPort:        8999,
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
