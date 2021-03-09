package dto

type ResponseProto struct {
	Result    interface{}    `sortId:"2" orderId:"1" generic:"true"`
	OutParam []interface{}	`sortId:"1" orderId:"2" generic:"true"`
}
