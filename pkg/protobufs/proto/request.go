package dto

import "container/list"

type RequestProto struct {
	Lookup    string     `sortId:"1" orderId:"1"`
	Method    string     `sortId:"2" orderId:"2"`
	ParamList *list.List `sortId:"3" orderId:"3"`
}
