

Server类包括名字Name,ip版本，监听的地址和端口，还包括
1. 消息处理器 MsgHandler
2. 连接管理器 ConnMgr
2. 连接创建前后的hook函数OnConnStart、OnConnStop


IRouter是一个抽象接口，包含处理IRequest的方法：
PreHandle  Handle  PostHandle


IRequest是一个抽象接口，包含一些方法：
获取当前连接：GetConnection
获取请求的消息数据: GetData
获取请求消息的ID: GetMsgID


Request类:
包含一个连接(conn)以及一个消息对象(msg)




MsgHandler类组件：
1.存放msgId和对应处理方法的Apis map[uint32]ziface.IRouter
2.存放消息的消息队列 TaskQueue []chan ziface.IRequest,这个切片中管道的数量就跟work协程的数量一致
3.work协程的数量 WorkerPoolSize




ConnManager类组件:
1.所有连接的集合，每个序号对应一个连接
connections map[uint32]ziface.IConnection
2.保护协程安全的读写锁



connection类：
TcpServer(记录当前连接属于哪个server)
Conn *net.TCPConn类型，表示当前tcp连接
还有连接id，连接状态，用于通讯的channel ExitChan chan bool,用于告诉写groutine连接已经关闭
msgChan chan []byte 用于读写协程之间的通信
还有个消息管理器
还有个连接属性集合  property map[string]interface{}  以及保护连接属性的锁
有启动、获取、停止连接，获取地址和端口，发送数据（SendMsg)
设置、获取、移除连接属性



Message类:
包含消息Id,消息长度(DataLen)，消息内容(Data[]byte)
消息Id和消息长度记录在包头，固定一共8个字节
此外还有一些列获取属性(长度，内容)和设置的方法





DataPack类：
包含获取包头长度以及封包、拆包(pack/unpack)的方法
封包是把IMessage对象转为二进制字节流，拆包反之





GlobalObj类，包含一些配置信息，通过Reload函数从json配置文件加载，如果json没有某个属性，我们还设置了默认值







项目执行准备工作：
首先调用NewServer初始化Server，通过全局配置文件指定Name,ip版本，监听的地址和端口，还New出了连接管理器和消息处理器。
然后注册hook函数，然后添加router,以便根据消息id走相应的router。定义了基类router:BaseRouter，派生出子类，(PingRouter和HelloRouter)子类可根据需要重写实现基类的接口。

项目执行流程：
执行Start,首先调用消息处理器的StartWorkerPool开启worker工作池,开启的协程数量由WorkerPoolSize决定，然后初始化TaskQueue，他是个切片，每个元素是装有IRequest的管道，管道数量也由WorkerPoolSize决定，每个管道最多可以容纳MaxWorkerTaskLen个数据。每个worker协程循环读取TaskQueue，注意每个协程度自己的TaskQueue。读到数据就调用DoMsgHandler处理。
然后创建一个TCP地址，之后监听这个地址。然后不断调用Accept,如果成功，判断是否达到最大连接数，然后new出一个连接对象，调用连接的启动（Start)方法启动连接，实际上是开启读写协程。


然后看开启读写协程之后的事情：
首先了解一下pack和unpack怎么实现，其中pack是封包，把收到的消息的长度先写入Buff缓冲区，然后写入MsgId，然后写入数据。最后以二进制文件的形式返回这个Buff。而unpack是创建一个Message对象，然后把二进制的消息长度和消息Id读入这个Message对象，unpack返回这个对象的时候，里面还没有实际的data数据

读协程开启之后，首先New出个拆包解包对象dp，然后创建一个字节流切片headData，根据HeadLen把读入的消息头放入headData。然后调用unpack(headData)返回一个msg消息对象，然后根据这个消息对象知道消息长度，接着把消息读入data里面(data []byte,msg.GetMsgLen())
然后直接调用msg.SetData(data)设置msg对象的data。此时创建一个Request的对象，用当前连接和当前消息(msg)构造。如果设置了WorkerPoolSize>0,那么就把这个Request放入TaskQueue，然后读到数据后找到对应路由处理请求。否则直接通过找到对应消息Id，然后通过对应路由处理这个请求。(handler,ok:=MsgHandle.Apis[request.GetMsgID()]),这个放入任务队列是通过取余放的，均衡放入，(mh.TaskQueue[workID]<--request),这个workID是%出来的。

然后看写协程：
一直监听msgChan，如果有数据，就发给客户端，同时也监听ExitChan，如果有数据，就退出。
什么时候msgChan有数据？当Request放入TaskQueue或者直接调用相应router处理的时候，handle处理函数会调用SendMsg把数据发给管道







