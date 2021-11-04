/*
 * @Author: your name
 * @Date: 2021-10-19 17:41:17
 * @LastEditTime: 2021-11-04 17:05:48
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/network/session/session.go
 */
package session

import "net"

type Session interface {
	Send(msg interface{}) error
	Start()
	Close()
	Conn() net.Conn
}
