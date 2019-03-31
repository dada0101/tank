package network

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const (
	BUF_LENG = 1024

	AGENT_STATE_ONLINE = 1
	AGENT_STATE_OUTLINE = 2
	AGENT_STATE_CLOSED = 3
)

func State(s int32) string {
	if s == AGENT_STATE_ONLINE {
		return "[AGENT_STATE_ONLINE]"
	} else if s == AGENT_STATE_CLOSED {
		return "[AGENT_STATE_CLOSED]"
	} else {
		return "[AGENT_STATE_CLOSED]"
	}
}

var (
	BUFEXCEED_ERR = errors.New("send buffer exceed!")
	CLOSED_AGENT_ERR = errors.New("the agent has closed!")

	CloseHandle func(a *Agent)
)

type (
	Agent struct {
		id int32		// 自己的id
		uid int32		// 玩家id
		noHBCnt uint32  // 无心跳计数，当到达某个值时，将连接断开
		conn net.Conn	// 连接
		lastMsgId uint32 // 最后一条消息id
		state int32		// 当前状态
		chDie chan struct{}
		sendBuf chan []byte		// buffer 缓冲
		ExtraChan chan struct{} // 用于两个agent简单同步
	}
)

func NewAgent(id int32, conn net.Conn) *Agent {
	return &Agent{
		id,
			0,
		0,
		conn,
		0,
		AGENT_STATE_ONLINE,
		make(chan struct{}),
		make(chan []byte, BUF_LENG),
		make(chan struct{}, 1),
	}
}

func (a *Agent) SetUid(uid int32) {
	a.uid = uid
}

func (a *Agent) Status()int32 {
	return atomic.LoadInt32(&a.state)
}
func (a *Agent)setStatus (s int32) {
	atomic.StoreInt32(&a.state, s)
}

func (a *Agent) IncNoHBCnt() {
	atomic.AddUint32(&a.noHBCnt, 1)
}
func (a *Agent) NoHBCntZero() {
	atomic.StoreUint32(&a.noHBCnt, 0)
}
func (a *Agent) GetHBCnt() uint32{
	return atomic.LoadUint32(&a.noHBCnt)
}

// no block send
func (a *Agent) Send(m []byte) error {
	if a.Status() == AGENT_STATE_CLOSED {
		return CLOSED_AGENT_ERR
	}
	if len(a.sendBuf) >= BUF_LENG {
		return BUFEXCEED_ERR
	}
	a.sendBuf <- m
	return nil
}


func (a *Agent) Close() error {
	if a.Status() == AGENT_STATE_CLOSED {
		return CLOSED_AGENT_ERR
	}
	a.setStatus(AGENT_STATE_CLOSED)
	select {
		case <- a.chDie:
	default:
	//	close(a.chDie)
	}
	// TODO: 这里设计问题，需要修改
	CloseHandle(a)
	return a.conn.Close()
}

func DumpAgent(a *Agent) (str string){
	str = "--Agent Dump--\n"
	str += fmt.Sprintf("id:\t%d\tlastHB:\t%d\n",a.id, a.noHBCnt)
	str += fmt.Sprintf("net-connect:\t%v\n", a.conn)
	str += fmt.Sprintf("state:%s\n", State(a.Status()))
	return str
}


func (a *Agent) Server(processing func(* Agent, []byte)([]byte, bool)) {
	buf := make([]byte, 2048)
	for {
		n, err := a.conn.Read(buf)
		if err != nil {
			log.Println("can not read buffer from this agent and agent will close. don't worry, it's not a bug :)")
			a.chDie <- struct{}{}
			break
		}
		if processing == nil {
			log.Println(buf[:n])
			continue
		}
		outBuf, isEmpty := processing(a, buf[:n])
		if !isEmpty {
			a.sendBuf <- outBuf
		}
	}
	log.Println(a.id, "\t agent server process finished")
}

func (a *Agent) Write(heartbeat time.Duration) {
	ticker := time.NewTicker(heartbeat)
	defer func() {
		ticker.Stop()
		close(a.sendBuf)
		a.Close()
		log.Println(a.id, "\t agent write process finished")
	}()
	for {
		select {
		 	case <-ticker.C:
		 		if a.GetHBCnt() > 4 {
		 			a.chDie <- struct{}{}
				}
				a.IncNoHBCnt()
		 	case data := <-a.sendBuf:
		 		_,err := a.conn.Write(data)
		 		if err != nil {
		 			log.Println(err)
		 			return
				}
		 	case <-a.chDie:
		 		return
		}
	}
}

func (a *Agent) Run(processing func(*Agent, []byte)([]byte, bool), timeout time.Duration) {
	go a.Server(processing)
	go a.Write(timeout)
}
