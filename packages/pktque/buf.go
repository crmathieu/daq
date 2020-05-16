package pktque

import (
	"github.com/crmathieu/daq/data"
)

const DataPointCount = 256

type Buf struct {
	Head, Tail BufPos
	pkts       []data.DataPoint
	Size       int
	Count      int
}

type BufPos int

func (b *Buf) Get(pos BufPos) data.DataPoint {
	return b.pkts[int(pos)&(len(b.pkts)-1)]
}

func (b *Buf) IsValidPos(pos BufPos) bool {
	return pos.GE(b.Head) && pos.LT(b.Tail)
}

func (bp BufPos) LT(pos BufPos) bool {
	return bp-pos < 0
}

func (bp BufPos) GE(pos BufPos) bool {
	return bp-pos >= 0
}

func (bp BufPos) GT(pos BufPos) bool {
	return bp-pos > 0
}

func NewBuf() *Buf {
	return &Buf{
		pkts: make([]data.DataPoint, DataPointCount),
	}
}

func (b *Buf) Pop() data.DataPoint { 
	if b.Count == 0 {
		panic("queue.Buf: Pop() when count == 0")
	}

	i := int(b.Head) & (len(b.pkts) - 1)
	dp := b.pkts[i]
	b.pkts[i] = data.DataPoint{} 
	b.Size -= data.DATAPOINT_SIZE 
	b.Head++
	b.Count--
	return dp
}

func (b *Buf) grow() {
	newpkts := make([]data.DataPoint, len(b.pkts)*2)
	for i := b.Head; i.LT(b.Tail); i++ {
		newpkts[int(i)&(len(newpkts)-1)] = b.pkts[int(i)&(len(b.pkts)-1)]
	}
	b.pkts = newpkts
}

func (b *Buf) Push(pkt data.DataPoint) {
	if b.Count == len(b.pkts) {
		b.grow()
	}
	b.pkts[int(b.Tail)&(len(b.pkts)-1)] = pkt
	b.Tail++
	b.Count++
	b.Size += data.DATAPOINT_SIZE
}

