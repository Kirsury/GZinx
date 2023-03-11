package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	//1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err: ", err)
		return
	}
	//2 从客户端读取数据，拆包处理

	go func() {
		//创建一个go承载 负责从客户端处理业务
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error")
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				//------> 拆包的过程 <---------
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1 第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2 第二次从conn读，根据head中的datalen 再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err : ", err)
							return
						}

						//完整的一个消息已经读取完毕
						fmt.Println("-----》 Recv MsgId: ", msg.Id, ", datalen = ", msg.DataLen, " data = ", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
	}

	//封装一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程， 封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
	}
	//封装第二个msg1包
	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error: ", err)
	}
	//将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	//一次性发送给服务端
	conn.Write(sendData1)

	//客户端阻塞
	select {}
}
