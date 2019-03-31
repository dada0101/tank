package proto

import (
	"reflect"
)

// handle connect message
// mothod	--- Heartbeat
//			--- Register
//			--- Login
//			--- Logout

// handle player message (battle)
// mothod	--- StartFight
//			--- UpdateUnitInfo
//			--- Shooting
//			--- Hit
//----------------------- (player)
//			--- GetScore
//			--- AddScore
// 			--- GetList
//			--- UpdateInfo
//			--- GetAchieve
//------------------------ (room)
//			--- GetRoomList
// 			--- CreateRoom
//			--- EnterRoom
//			--- GetRoomInfo
//			--- LeaveRoom
//			--- Prepare
//			--- Cancel
//			--- Chat
//			--- ChatName

type (

	MethodDescriptor struct {
		name string					// 方法名
		paramsType []reflect.Kind	// 参数类型
		method interface{}			// 方法
	}
	NameMapping struct {
		mapping map[string]*MethodDescriptor
	}
)

var (
	rpcMap NameMapping
)

func GetGRpcMap() *NameMapping{
	return &rpcMap
}
/*
func Init() {
	GRpcMap = & NameMapping{
		make(map[string]*MethodDescriptor),
	}
}
*/

func NewMethodDescriptor() *MethodDescriptor{
	return &MethodDescriptor{
		"",
		nil,
		nil,
	}
}
func(md *MethodDescriptor) SetMethod(name string,  pt []reflect.Kind , method interface{}) {
	md.name = name
	md.method = method
	md.paramsType = pt
}

func(md *MethodDescriptor) NumIn() int {
	return len(md.paramsType)
}

func (md *MethodDescriptor)In(i int) reflect.Kind {
	return md.paramsType[i]
}

func (md *MethodDescriptor)GetMethod() interface{}{
	return md.method
}

func (rpcMap *NameMapping)Init() {
	rpcMap.mapping = make(map[string]*MethodDescriptor)
}

func(rpcMap *NameMapping) AddActionHandle(name string, pt []reflect.Kind, ah interface{}) {
	md := NewMethodDescriptor()
	md.SetMethod(name, pt, ah)
	rpcMap.mapping[name] = md
}

func(rpcMap *NameMapping) GetMethodDescriptor(name string) (*MethodDescriptor, bool) {
	md, found := rpcMap.mapping[name]
	return md , found
}