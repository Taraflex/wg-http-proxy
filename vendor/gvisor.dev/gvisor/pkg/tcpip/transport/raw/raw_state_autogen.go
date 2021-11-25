// automatically generated by stateify.

package raw

import (
	"gvisor.dev/gvisor/pkg/state"
	"gvisor.dev/gvisor/pkg/tcpip/buffer"
)

func (p *rawPacket) StateTypeName() string {
	return "pkg/tcpip/transport/raw.rawPacket"
}

func (p *rawPacket) StateFields() []string {
	return []string{
		"rawPacketEntry",
		"data",
		"receivedAt",
		"senderAddr",
		"packetInfo",
	}
}

func (p *rawPacket) beforeSave() {}

// +checklocksignore
func (p *rawPacket) StateSave(stateSinkObject state.Sink) {
	p.beforeSave()
	var dataValue buffer.VectorisedView
	dataValue = p.saveData()
	stateSinkObject.SaveValue(1, dataValue)
	var receivedAtValue int64
	receivedAtValue = p.saveReceivedAt()
	stateSinkObject.SaveValue(2, receivedAtValue)
	stateSinkObject.Save(0, &p.rawPacketEntry)
	stateSinkObject.Save(3, &p.senderAddr)
	stateSinkObject.Save(4, &p.packetInfo)
}

func (p *rawPacket) afterLoad() {}

// +checklocksignore
func (p *rawPacket) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &p.rawPacketEntry)
	stateSourceObject.Load(3, &p.senderAddr)
	stateSourceObject.Load(4, &p.packetInfo)
	stateSourceObject.LoadValue(1, new(buffer.VectorisedView), func(y interface{}) { p.loadData(y.(buffer.VectorisedView)) })
	stateSourceObject.LoadValue(2, new(int64), func(y interface{}) { p.loadReceivedAt(y.(int64)) })
}

func (e *endpoint) StateTypeName() string {
	return "pkg/tcpip/transport/raw.endpoint"
}

func (e *endpoint) StateFields() []string {
	return []string{
		"DefaultSocketOptionsHandler",
		"transProto",
		"waiterQueue",
		"associated",
		"net",
		"stats",
		"ops",
		"rcvList",
		"rcvBufSize",
		"rcvClosed",
		"frozen",
	}
}

// +checklocksignore
func (e *endpoint) StateSave(stateSinkObject state.Sink) {
	e.beforeSave()
	stateSinkObject.Save(0, &e.DefaultSocketOptionsHandler)
	stateSinkObject.Save(1, &e.transProto)
	stateSinkObject.Save(2, &e.waiterQueue)
	stateSinkObject.Save(3, &e.associated)
	stateSinkObject.Save(4, &e.net)
	stateSinkObject.Save(5, &e.stats)
	stateSinkObject.Save(6, &e.ops)
	stateSinkObject.Save(7, &e.rcvList)
	stateSinkObject.Save(8, &e.rcvBufSize)
	stateSinkObject.Save(9, &e.rcvClosed)
	stateSinkObject.Save(10, &e.frozen)
}

// +checklocksignore
func (e *endpoint) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &e.DefaultSocketOptionsHandler)
	stateSourceObject.Load(1, &e.transProto)
	stateSourceObject.Load(2, &e.waiterQueue)
	stateSourceObject.Load(3, &e.associated)
	stateSourceObject.Load(4, &e.net)
	stateSourceObject.Load(5, &e.stats)
	stateSourceObject.Load(6, &e.ops)
	stateSourceObject.Load(7, &e.rcvList)
	stateSourceObject.Load(8, &e.rcvBufSize)
	stateSourceObject.Load(9, &e.rcvClosed)
	stateSourceObject.Load(10, &e.frozen)
	stateSourceObject.AfterLoad(e.afterLoad)
}

func (l *rawPacketList) StateTypeName() string {
	return "pkg/tcpip/transport/raw.rawPacketList"
}

func (l *rawPacketList) StateFields() []string {
	return []string{
		"head",
		"tail",
	}
}

func (l *rawPacketList) beforeSave() {}

// +checklocksignore
func (l *rawPacketList) StateSave(stateSinkObject state.Sink) {
	l.beforeSave()
	stateSinkObject.Save(0, &l.head)
	stateSinkObject.Save(1, &l.tail)
}

func (l *rawPacketList) afterLoad() {}

// +checklocksignore
func (l *rawPacketList) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &l.head)
	stateSourceObject.Load(1, &l.tail)
}

func (e *rawPacketEntry) StateTypeName() string {
	return "pkg/tcpip/transport/raw.rawPacketEntry"
}

func (e *rawPacketEntry) StateFields() []string {
	return []string{
		"next",
		"prev",
	}
}

func (e *rawPacketEntry) beforeSave() {}

// +checklocksignore
func (e *rawPacketEntry) StateSave(stateSinkObject state.Sink) {
	e.beforeSave()
	stateSinkObject.Save(0, &e.next)
	stateSinkObject.Save(1, &e.prev)
}

func (e *rawPacketEntry) afterLoad() {}

// +checklocksignore
func (e *rawPacketEntry) StateLoad(stateSourceObject state.Source) {
	stateSourceObject.Load(0, &e.next)
	stateSourceObject.Load(1, &e.prev)
}

func init() {
	state.Register((*rawPacket)(nil))
	state.Register((*endpoint)(nil))
	state.Register((*rawPacketList)(nil))
	state.Register((*rawPacketEntry)(nil))
}
