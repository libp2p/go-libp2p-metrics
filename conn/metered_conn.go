package meterconn

import (
	metrics "github.com/libp2p/go-libp2p-metrics"
	tpt "github.com/libp2p/go-libp2p-transport"
)

func WrapConn(bwc metrics.Reporter, conn tpt.Conn) tpt.Conn {
	switch c := conn.(type) {
	case tpt.SingleStreamConn:
		return newMeteredSingleStreamConn(c, bwc.LogRecvMessage, bwc.LogSentMessage)
	case tpt.MultiStreamConn:
		return newMeteredMultiStreamConn(c, bwc.LogRecvMessage, bwc.LogSentMessage)
	default:
		panic("c is neither a SingleStreamConn nor a MultiStreamConn")
	}
}
