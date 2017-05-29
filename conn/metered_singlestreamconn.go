package meterconn

import (
	metrics "github.com/libp2p/go-libp2p-metrics"
	tpt "github.com/libp2p/go-libp2p-transport"
)

type MeteredSingleStreamConn struct {
	tpt.SingleStreamConn

	mesRecv metrics.MeterCallback
	mesSent metrics.MeterCallback
}

var _ tpt.SingleStreamConn = &MeteredSingleStreamConn{}

func newMeteredSingleStreamConn(base tpt.SingleStreamConn, rcb metrics.MeterCallback, scb metrics.MeterCallback) *MeteredSingleStreamConn {
	return &MeteredSingleStreamConn{
		SingleStreamConn: base,
		mesRecv:          rcb,
		mesSent:          scb,
	}
}

func (mc *MeteredSingleStreamConn) Read(b []byte) (int, error) {
	n, err := mc.SingleStreamConn.Read(b)
	mc.mesRecv(int64(n))
	return n, err
}

func (mc *MeteredSingleStreamConn) Write(b []byte) (int, error) {
	n, err := mc.SingleStreamConn.Write(b)
	mc.mesSent(int64(n))
	return n, err
}
