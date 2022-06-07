package model

//用于不同线程之间通信(客户端在主线程发送请求后，接收其他goroutine收到的回复),
//注意发送后必须有且只有一个goroutine正在接收，否者可能引发未知错误
var Ch1 chan any = make(chan any)