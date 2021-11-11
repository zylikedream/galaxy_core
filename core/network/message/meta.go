package message

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

var (
	// 消息元信息与消息名称，消息ID和消息类型的关联关系
	metaByFullName = map[string]*MessageMeta{}
	metaByID       = map[int]*MessageMeta{}
	metaByType     = map[reflect.Type]*MessageMeta{}
)

type MessageMeta struct {
	ID       int
	FullName string
	Type     reflect.Type
}

func fullName(t reflect.Type) string {

	var sb strings.Builder
	sb.WriteString(path.Base(t.PkgPath()))
	sb.WriteString(".")
	sb.WriteString(t.Name())

	return sb.String()
}

func (m *MessageMeta) TypeName() string {

	if m == nil {
		return ""
	}

	return m.Type.Name()
}

func (m *MessageMeta) NewInstance() interface{} {
	if m.Type == nil {
		return nil
	}

	return reflect.New(m.Type).Interface()
}

func RegisterMessageMeta(ID int, msg interface{}) *MessageMeta {
	meta := &MessageMeta{
		ID:   ID,
		Type: reflect.TypeOf(msg),
	}
	// 注册时, 统一为非指针类型
	if meta.Type.Kind() == reflect.Ptr {
		meta.Type = meta.Type.Elem()
	}
	meta.FullName = fullName(meta.Type)

	if _, ok := metaByType[meta.Type]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by type: %d name: %s", meta.ID, meta.Type.Name()))
	} else {
		metaByType[meta.Type] = meta
	}

	if _, ok := metaByFullName[meta.FullName]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by fullname: %s", meta.FullName))
	} else {
		metaByFullName[meta.FullName] = meta
	}

	if meta.ID == 0 {
		panic("message meta require 'ID' field: " + meta.TypeName())
	}

	if prev, ok := metaByID[meta.ID]; ok {
		panic(fmt.Sprintf("Duplicate message meta register by id: %d type: %s, pre type: %s", meta.ID, meta.TypeName(), prev.TypeName()))
	} else {
		metaByID[meta.ID] = meta
	}

	return meta
}

func MessageMetaByID(id int) *MessageMeta {
	if v, ok := metaByID[id]; ok {
		return v
	}

	return nil
}

func MessageMetaByType(t reflect.Type) *MessageMeta {

	if t == nil {
		return nil
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if v, ok := metaByType[t]; ok {
		return v
	}

	return nil
}

// 根据消息对象获得消息元信息
func MessageMetaByMsg(msg interface{}) *MessageMeta {

	if msg == nil {
		return nil
	}

	return MessageMetaByType(reflect.TypeOf(msg))
}
