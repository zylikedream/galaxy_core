/*
 * @Author: your name
 * @Date: 2021-10-19 17:41:17
 * @LastEditTime: 2021-11-04 17:05:48
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/gxynet/conn/conn.go
 */
package conn

import (
	"context"
	"net"
)

type Conn interface {
	Send(msg interface{}) error
	Start(ctx context.Context) error
	Close(context.Context, error)
	Raw() net.Conn
	GetData() interface{}
	SetData(interface{})
}
