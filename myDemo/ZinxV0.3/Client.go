package main

import (
	"fmt"
	"net"
	"time"
)

/*
模拟客户端
*/

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	// 1直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		// 2连接调用write，写数据
		_, err := conn.Write([]byte("hello zinx v0.3"))
		if err != nil {
			fmt.Println("write conn err: ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		fmt.Printf("serever call back: %s, cnt = %d\n", buf, cnt)

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
