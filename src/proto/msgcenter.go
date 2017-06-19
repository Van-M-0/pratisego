package proto

import (
	"reflect"
	"fmt"
)

type MessageCenter struct {
	types 			map[uint32]reflect.Type
	ids 			map[reflect.Type]uint32
}

func NewMessageCenter() *MessageCenter {
	return &MessageCenter{
		types: 		make(map[uint32]reflect.Type),
		ids: 		make(map[reflect.Type]uint32),
	}
}

func (m *MessageCenter) Register(id uint32, msg interface{}) {
	rt := reflect.TypeOf(msg)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	m.types[id] = rt
	m.ids[rt] = id
}

func (m *MessageCenter) GetType(id uint32) (reflect.Type, error) {
	if rt, ok := m.types[id]; ok {
		return rt, nil
	} else {
		return nil, fmt.Errorf("invalid msg id %d", id)
	}
}

func (m *MessageCenter) NewMessage(id uint32) (interface{}, error) {
	if t, ok := m.types[id]	; ok {
		return reflect.New(t).Interface(), nil
	} else {
		return nil, fmt.Errorf("invalid msg id %d", id)
	}
}
