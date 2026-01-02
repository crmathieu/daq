package streamer

import (
	"github.com/crmathieu/daq/packages/data"
	"fmt"
)

const DataPointCount = 256
type quBuf struct {
	Head, Tail BufIdx
//	pkts       []data.DataPoint
	pkts       []*[data.PACKET_GRP]data.DataPoint
	Size       int
	Count      int
}

type BufIdx int

func (b *quBuf) Get(pos BufIdx) [data.PACKET_GRP]data.DataPoint {
//	return b.pkts[int(pos)&(len(b.pkts)-1)]
	return *b.pkts[int(pos)&(len(b.pkts)-1)]
}

func (b *quBuf) IsValidPos(pos BufIdx) bool {
	return pos.GE(b.Head) && pos.LT(b.Tail)
}

func (bp BufIdx) LT(pos BufIdx) bool {
	return bp-pos < 0
}

func (bp BufIdx) GE(pos BufIdx) bool {
	return bp-pos >= 0
}

func (bp BufIdx) GT(pos BufIdx) bool {
	return bp-pos > 0
}

func NewBuf() *quBuf {
	pktsp := make([]*[data.PACKET_GRP]data.DataPoint, DataPointCount)
	pktsdata := make([][data.PACKET_GRP]data.DataPoint, DataPointCount)
	for k := 0; k<DataPointCount; k++ {
//		pktsp[k] = &(make([]data.DataPoint, data.PACKET_GRP)) //&pktsdata[k]
		//p := unsafe.Pointer(&pktsdata[k]) //&(make([]data.DataPoint, data.PACKET_GRP)))
		pktsp[k] = &pktsdata[k]// (*[data.PACKET_GRP]data.DataPoint)(p)
		//pktsp[k] = &([data.PACKET_GRP]data.DataPoint)unsafe.Pointer(&(make([]data.DataPoint, data.PACKET_GRP))) //&pktsdata[k]
	}
	return &quBuf{
//		pkts: make([]data.DataPoint, DataPointCount),
		pkts: pktsp, //make([]*[data.PACKET_GRP]data.DataPoint, DataPointCount),
	}
}

func (b *quBuf) Pop() [data.PACKET_GRP]data.DataPoint { 
	if b.Count == 0 {
		panic("quBuf: Pop() when count == 0")
	}

	i := int(b.Head) & (len(b.pkts) - 1)
	dp := b.pkts[i]
//	b.pkts[i] = data.DataPoint{} 
	b.Size -= data.DATAPOINT_SIZE 
	b.Head++
	b.Count--
//	return dp
	return *dp
}

func (b *quBuf) grow() {
//	newpkts := make([]data.DataPoint, len(b.pkts)*2)
	fmt.Println("GROW!")
	newpkts := make([]*[data.PACKET_GRP]data.DataPoint, len(b.pkts)*2)
	newpktsdata := make([][data.PACKET_GRP]data.DataPoint, len(b.pkts)*2)

	for k := 0; k<len(b.pkts)*2; k++ {
		newpkts[k] = &newpktsdata[k]// (*[data.PACKET_GRP]data.DataPoint)(p)
	}
	for i := b.Head; i.LT(b.Tail); i++ {
//		newpkts[int(i)&(len(newpkts)-1)] = b.pkts[int(i)&(len(b.pkts)-1)]

		newpkts[int(i)&(len(newpkts)-1)] = b.pkts[int(i)&(len(b.pkts)-1)]
	}
	b.pkts = newpkts
}

//func (b *quBuf) Push(pkt data.DataPoint) {
func (b *quBuf) Push(pkt *[data.PACKET_GRP]data.DataPoint) {
	if b.Count == len(b.pkts) {
		b.grow()
	}
//	fmt.Println("PUSH len:", len(*pkt), "- expected:", len(*(b.pkts[int(b.Tail)&(len(b.pkts)-1)])))
//	b.pkts[int(b.Tail)&(len(b.pkts)-1)] = pkt
	*(b.pkts[int(b.Tail)&(len(b.pkts)-1)]) = *pkt
	b.Tail++
	b.Count++
	b.Size += data.DATAPOINT_SIZE
}

