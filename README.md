# rhino
go网络基础框架
# 使用领域
rhino励志打造一个完善的游戏服务器框架，其核心内容在于
· actor 主要服务游戏逻辑部分,actor内部是单线程，例如：日志服务，db服务，游戏服务，网关等等
· network 用于tcp，http，udp的通信，封装读写包
· 其他部分是封装了一些重复劳动的工具

# acotr模型
1 核心api


# network
服务端: 通过设置SetReadTimeout来获得心跳，而无须在上层做timer，可以通过test查看例子

客户端：保持默认设置即可

network.stream封装了len+body的读取方法，由于跨协程read不安全所以socket一般不会直接面向应用层逻辑，借助actor模型 由actor对其控制，来更好的解决和应对分布式中环路堵塞和短暂丢包的处理



# 编写一个简单的服务器
```go
  network.StartTcpServer(":8088", network.OptionHandler(func(conn net.Conn) (err error) {
		Stage().ActorOf(WithRemoteStream(func(ctx ActorContext) {
			switch body := ctx.Any().(type) {
			case *Started:
				fmt.Println("connect addr:", conn.RemoteAddr().String())
			case *Stopped:
			case []byte:
				fmt.Println("message: ", network.ReadBegin(body))
				ctx.Send(ctx.Self(), network.WriteBegin(0x102).Flush())
			case error: //设置了心跳会被通知
			default:
				fmt.Printf("untreated type %T \n", body)
			}
		}, conn))
		return
	}))

	cli := Stage().ActorOf(WithRemoteAddr(func(ctx ActorContext) {
		switch body := ctx.Any().(type) {
		case *Started:
			fmt.Println("open ok")
		case *Stopped:
		case []byte:
			fmt.Println("cli respond: ", network.ReadBegin(body))
		case Failure: //最终错误退出的原因

		default:
			fmt.Printf("other object miss handle %T \n", body)
		}
	}, "localhost:8088"))

	psend := network.WriteBegin(0x101, "who's your daddy!")
	cli.Tell(psend.Flush())
```
