package streamer

import (
	"github.com/crmathieu/daq/packages/data"
)

const DataPointCount = 256

type quBuf struct {
	Head, Tail BufIdx
	pkts       []data.DataPoint
	Size       int
	Count      int
}

type BufIdx int

func (b *quBuf) Get(pos BufIdx) data.DataPoint {
	return b.pkts[int(pos)&(len(b.pkts)-1)]
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
	return &quBuf{
		pkts: make([]data.DataPoint, DataPointCount),
	}
}

func (b *quBuf) Pop() data.DataPoint { 
	if b.Count == 0 {
		panic("quBuf: Pop() when count == 0")
	}

	i := int(b.Head) & (len(b.pkts) - 1)
	dp := b.pkts[i]
	b.pkts[i] = data.DataPoint{} 
	b.Size -= data.DATAPOINT_SIZE 
	b.Head++
	b.Count--
	return dp
}

func (b *quBuf) grow() {
	newpkts := make([]data.DataPoint, len(b.pkts)*2)
	for i := b.Head; i.LT(b.Tail); i++ {
		newpkts[int(i)&(len(newpkts)-1)] = b.pkts[int(i)&(len(b.pkts)-1)]
	}
	b.pkts = newpkts
}

func (b *quBuf) Push(pkt data.DataPoint) {
	if b.Count == len(b.pkts) {
		b.grow()
	}
	b.pkts[int(b.Tail)&(len(b.pkts)-1)] = pkt
	b.Tail++
	b.Count++
	b.Size += data.DATAPOINT_SIZE
}

