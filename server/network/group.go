package network

import (
	"log"
	"sync"
	"sync/atomic"
)

const (
	GROUP_STATUS_WORKING = 0
	GROUP_STATUS_CLOSED = 1

)


type Group struct {
	mu sync.Mutex
	id uint32		// group id
	status int32	// group 状态
	agents map[*Agent]bool // 组中所有的agent
}

func(group *Group)isClosed() bool {
	return atomic.LoadInt32(&group.status) == GROUP_STATUS_CLOSED
}
func(group *Group)setClosed() {
	atomic.StoreInt32(&group.status, GROUP_STATUS_CLOSED)
}


func NewGroup(id uint32)*Group {
	return &Group{
		id: id,
		status: GROUP_STATUS_WORKING,
		agents: make(map[*Agent]bool),
	}
}

func(group *Group)Members() []*Agent {
	members := make([]*Agent, 0)
	for a, _ := range group.agents {
		members = append(members, a)
	}
	return members
}

func(group *Group)Add(agent *Agent) bool{
	if group.isClosed() {
		log.Println("this group has closed, but someone want to add agent to it.")
		return false
	}
	group.mu.Lock()
	defer group.mu.Unlock()
	group.agents[agent] = true
	return true
}

func (group* Group)Del(agent *Agent) bool{
	if group.isClosed() {
		log.Println("this group has closed, but someone want to delete agent from it.")
		return false
	}
	group.mu.Lock()
	defer group.mu.Unlock()
	delete(group.agents, agent)
	return true
}

func(group *Group)Broadcast(msg []byte) {
	if group.isClosed() {
		log.Println("this group has closed, don't broadcast message.")
		return
	}
	group.mu.Lock()
	defer group.mu.Unlock()
	for a, _ := range group.agents {
		if a.Status() == AGENT_STATE_CLOSED {
			continue
		}
		// TODO 可能有问题
		if err := a.Send(msg); err != nil {
			log.Println(err)
		}
	}
}


func (group *Group) Close() {
	if group.isClosed() {
		log.Println("this group has closed, but someone want to close it again.")
		return
	}
	group.setClosed()
	group.agents = nil
}