package server

import (
	"bufio"
	"encoding/json"
	"io"
	"kde/config"
	"log"
	"net"
	"strings"
	"time"
)

// 设备
type Terminal struct {
	lastTimeStamp int64 //最后时间戳
	flowId        int64
	data          string   //数据
	macType       string   //设备类型
	sn            string   //设备编号
	Conn          net.Conn //连接
}

//web显示设备信息
type Device struct {
	Sn            string `json:"sn"`
	LastTimeStamp string `json:"lastTimeStamp"`
	FlowId        int64  `json:"flowId"`
	Data          string `json:"data"`
	MacType       string `json:"macType"`
}

//推送到mqtt的消息
type DataPo struct {
	Time string `json:"time"`
	Src  string `json:"src"`
}

//连接管理
var connManager map[string]*Terminal

//处理消息
func recvConnMsg(conn net.Conn) {
	addr := conn.RemoteAddr()
	var term = &Terminal{
		Conn:          conn,
		lastTimeStamp: time.Now().Unix() * 1000,
	}
	term.Conn = conn
	connManager[addr.String()] = term

	defer func() {
		delete(connManager, addr.String())
		conn.Close()
		//todo 维护连接列表
		for k, v := range connManager {
			log.Println(k, v)
		}
	}()
	for {
		// ReadString 会一直阻塞直到遇到分隔符 '\n'
		// 遇到分隔符后会返回上次遇到分隔符或连接建立后收到的所有数据, 包括分隔符本身
		// 若在遇到分隔符之前遇到异常, ReadString 会返回已收到的数据和错误信息
		//超过1024连接会关闭
		msg, err := bufio.NewReader(&io.LimitedReader{R: conn, N: 1024}).ReadString('\n')
		if err != nil {
			// 通常遇到的错误是连接中断或被关闭，用io.EOF表示
			if err == io.EOF {
				log.Println("connection close")
			} else {
				log.Println(err)
			}
			return
		}
		log.Printf("ip:%s,rcv:%s\n", addr.String(), msg)
		//验证数据
		//先找一下P*和*B*K，并且P*在*B*K之前
		p := strings.Index(msg, config.BeginStr)
		k := strings.Index(msg, config.EndStr)
		if p == -1 || k == -1 || k < p {
			log.Println("非法PBK:", msg)
			return
		}
		//截取P...K 以P*开头和*B*K结尾的才是正确协议 并且只能以单个*号分割
		msg = msg[p : k+4]
		strS := strings.Split(msg, "*")
		for _, v := range strS {
			if v == "" {
				log.Println("非法**:", msg)
				return
			}
		}
		term.sn = strS[4]
		term.macType = strS[3]
		//架桥机46，履带吊46，设备定位45，龙门吊48 长度限定10-100
		l := len(strS)
		if l < config.MinLength || l > config.MaxLength {
			//log.Println("长度:", l)
			log.Println("非法len:", msg)
			return
		}
		//验证设备列表 （取redis里面数据）
		val, err := config.GVA_REDIS.Get(config.Ctx, "macList").Result()
		if err != nil {
			log.Println(err)
		}
		//log.Println(val)
		var macList []string
		err = json.Unmarshal([]byte(val), &macList)
		if err != nil {
			log.Println(err)
		}
		//log.Println(macList)
		//log.Println(len(macList))
		if len(macList) == 0 {
			log.Println("非法设备list为空:", msg)
			return
		}
		//防止数据太频繁
		now := time.Now().Unix() * 1000
		lastTime := term.lastTimeStamp
		//log.Println(now - lastTime)
		if now-lastTime < config.MIN_TIME && term.flowId > 0 {
			log.Println("非法频率:", msg)
			continue
		}
		// 将收到的信息发送给客户端
		conn.Write([]byte(msg))
		dataPo := DataPo{
			Time: time.Now().Format("2006-01-02 15:04:05"),
			Src:  msg,
		}
		jsons, errs := json.Marshal(dataPo) //转换成JSON返回的是byte[]
		if errs != nil {
			log.Println(errs.Error())
		}
		//第一次连接，data有变化，data无变化超过5分钟 这三种情况下推送消息
		if term.flowId == 0 || term.data != msg || now-lastTime > config.MAX_TIME {
			term.lastTimeStamp = time.Now().Unix() * 1000
			term.data = msg
			term.flowId++
			//发送至mqtt
			config.GVA_MQTT.Publish(strS[4], 1, false, string(jsons))
		}
	}
}

//建立连接
func TCPServer(addr string) {
	connManager = make(map[string]*Terminal)

	listenSock, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}
	defer listenSock.Close()

	for {
		newConn, err := listenSock.Accept()
		if err != nil {
			continue
		}
		go recvConnMsg(newConn)
	}
}
