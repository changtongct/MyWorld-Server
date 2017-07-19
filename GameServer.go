package main

import (
	"fmt"
	"os"
	"net"
//	"unsafe"
	"strconv"
	"github.com/go-redis/redis"
)

//func connectRedis() *redis.Client {
//func addEntity(client *redis.Client, newentity EntityInfo) {
//func getEntity(client *redis.Client, key string) map[string]interface{} {

var new_id int32 = 666
var bullet_id int32 = 5000

//错误检查
func recordError(err error, info string) {
	fmt.Println(info + " " + err.Error())
}

func logIn(conn net.Conn, messages chan InternetPackage, client *redis.Client) int32 {
	//设置新id
	m_id := new_id
	//发送生成的id
	package_id := EntityToBytes(
					InternetPackage{check: CHECK_BYTE, ptype: 1,
					state: LOG_ALLOC_ID, reserve: 0, id: m_id,
					X: 0, Y: 0, Z: 0, toX: 0, toY: 0, toZ: 0})
	_, err := conn.Write(package_id)
	if err != nil {
		return -1
	}

	//接收id确认
	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		return -2
	}
	if buf[0] != CHECK_BYTE {
		return -3
	}
	//校验客户端是否收到分配id
	if BytesToInt32(buf[4:8]) != m_id {
		return -4
	}
	package_login := InternetPackage{check: CHECK_BYTE, ptype: 1,
                            state: LOG_IN, reserve: 0, id: m_id,
                            X: BytesToFloat32(buf[8:12]),
							Y: BytesToFloat32(buf[12:16]),
							Z: BytesToFloat32(buf[16:20]),
							toX: BytesToFloat32(buf[20:24]),
							toY: BytesToFloat32(buf[24:28]),
							toZ: BytesToFloat32(buf[28:32]),}
	//向所有用户广播上线信息
	messages <- package_login

	keys, _ := RedisGetAllEntityKeys(client)
	for _, v := range keys {
		entityinfo, _ := RedisGetEntity(client, v)
		tempid, _ := strconv.Atoi(v)
		tempX, _  := strconv.ParseFloat(entityinfo["X"], 64)
		tempY, _  := strconv.ParseFloat(entityinfo["Y"], 64)
		tempZ, _  := strconv.ParseFloat(entityinfo["Z"], 64)
		temptoX, _ := strconv.ParseFloat(entityinfo["toX"], 64)
		temptoY, _ := strconv.ParseFloat(entityinfo["toY"], 64)
		temptoZ, _ := strconv.ParseFloat(entityinfo["toZ"], 64)
		package_exists := InternetPackage{check: CHECK_BYTE, ptype: 1,
                            state: LOG_IN, reserve: 0, id: int32(tempid),
                            X: float32(tempX), Y: float32(tempY), Z: float32(tempZ),
							toX: float32(temptoX), toY: float32(temptoY), toZ: float32(temptoZ)}
//		fmt.Printf("EXISTS->%d,%f,%f,%f\n",package_exists.id,package_exists.X,package_exists.Y,package_exists.Z)
		p := EntityToBytes(package_exists)
		_, err := conn.Write(p)
		if err != nil {
			return -5
		}
	}
	//用户信息写入数据库
	RedisAddEntity(client, package_login)
	new_id += 1
	return m_id
}

func logOff(m_id int32, messages chan InternetPackage, client *redis.Client) {
	package_logoff := InternetPackage{check: CHECK_BYTE, ptype: 1,
                            state: LOG_OFF, reserve: 0, id: m_id,
                            X: 0, Y: 0, Z: 0, toX: 0, toY: 0, toZ: 0}
	messages <- package_logoff
	RedisDeleteEntity(client, strconv.Itoa(int(m_id)))
}

//服务器端接收数据线程
//参数: 数据连接conn 通讯通道messages
func Handler(conn net.Conn, messages chan InternetPackage, bullet_chan chan InternetPackage) int  {
	session_client,_ := RedisConnect()

	//打印新连接信息
	fmt.Println("connect from -->", conn.RemoteAddr().String())
	defer fmt.Println(conn.RemoteAddr().String(), "--> closed")

	m_id := logIn(conn, messages, session_client)
	if m_id < 0 {
		fmt.Println("登陆失败,返回值：", m_id)
		conn.Close()
		return -1
	}
	fmt.Println("登陆成功")
	////
	defer logOff(m_id, messages, session_client)

	buf := make([]byte, 1024)//设的比较大，防止溢出
	for {
		lengh, err := conn.Read(buf)
		if err != nil {
			conn.Close()
			return -1
		}
		if lengh > 0 {
			buf[lengh] = 0
		}
		if buf[0] != CHECK_BYTE {
			fmt.Println("非法连接,正在断开...")
			conn.Close()
			return -1
		}

		p := InternetPackage {	check: buf[0],
								ptype: buf[1],
								state: buf[2],
								reserve:buf[3],
								id: BytesToInt32(buf[4:8]),
								X: BytesToFloat32(buf[8:12]),
								Y: BytesToFloat32(buf[12:16]),
								Z: BytesToFloat32(buf[16:20]),
								toX: BytesToFloat32(buf[20:24]),
								toY: BytesToFloat32(buf[24:28]),
								toZ: BytesToFloat32(buf[28:32])}

		switch p.ptype {
			case 1:
				RedisChangeEntity(session_client, p)
				messages <- p
			case 2:
				bullet_chan <-p
//			case 3:
//				getEntity()
		}

//		mymap := getEntity(session_client, "234")
//		fmt.Println(mymap)
//		receiveStr := string(buf[0:lengh])
//		fmt.Println(receiveStr[0],receiveStr[1])

//		sendStr := "I received !!!"
//		conn.Write([]byte(sendStr))
	}
	return 0
}

func bulletManagement(conns *map[string]net.Conn, bullet_chan chan InternetPackage) {
	for {
		bullet_package := <-bullet_chan
		bullet_package.id = bullet_id
		bullet_package.state = 30
		bullet_id += 1
		if bullet_id > 5999 {
			bullet_id = 5000
		}
		s_bullet_package := EntityToBytes(bullet_package)
		for key, value := range *conns {
			_, err := value.Write(s_bullet_package)
			if err != nil {
				delete(*conns,key)
			}
		}
	}
}

func broadcastPlayerCoor(conns *map[string]net.Conn, messages chan InternetPackage) {
	for {
		p := <-messages
//		fmt.Printf("%x,%x,%x,%x,%d,%f,%f,%f,%f,%f,%f\n",p.check,p.ptype,p.state,p.reserve,p.id, p.X,p.Y,p.Z,p.toX,p.toY,p.toZ)
		s_p := EntityToBytes(p)
//		fmt.Println(s_p)
		for key, value := range *conns {
			//打印即将转发对象
//			fmt.Println("即将转发给->", key)
			_, err := value.Write(s_p)
			if err != nil {
				recordError(err, "Write")
				delete(*conns,key)
			}
		}
	}
}

func StartServer(port string) {
	tcpaddr := ":" + port
	//获取server套接字
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpaddr)
	if err != nil {
		recordError(err, "ResolveTCPAddr")
	}

	//开启监听
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		recordError(err, "ListenTCP")
	}

	messages := make(chan InternetPackage, 10)
	bullet_chan := make(chan InternetPackage, 30)
	//连接记录
	conns := make(map[string]net.Conn)
	go broadcastPlayerCoor(&conns, messages)
	go bulletManagement(&conns, bullet_chan)

	for {
		fmt.Println("监听中...")
		conn, err := l.Accept()
		if err != nil {
			recordError(err, "Accept")
		}

		conns[conn.RemoteAddr().String()] = conn

		//启动会话线程
		go Handler(conn, messages, bullet_chan)
	}
}

//参数: [port]
func main() {
	if len(os.Args) != 2 {
		fmt.Println("端口错误")
		os.Exit(0)
	}
	StartServer(os.Args[1])
}


