// Package pubsub implements publisher-subscribers model used in multi-channel streaming.
package streamer

import (
	"github.com/crmathieu/daq/packages/data"	
	"io"
	"sync"
)

//        time
// ----------------->
//
// V-A-V-V-A-V-V-A-V-V
// |                 |
// 0        5        10
// head             tail
// oldest          latest
//

// One publisher and multiple subscribers thread-safe packet buffer queue.
type Queue struct {
	buf                      *quBuf
	head, tail               int
	lock                     *sync.RWMutex
	cond                     *sync.Cond
	dpidx                 	 int
	closed                   bool
}

type QueueCursor struct {
	que    *Queue
	pos    BufPos
	gotpos bool
	init   func(buf *quBuf, dpidx int) BufPos
}

func (que *Queue) newCursor() *QueueCursor {
	return &QueueCursor{
		que: que,
	}
}

// NewQueue -------------------------------------------------------------------
// returns a new queue
// ----------------------------------------------------------------------------
func NewQueue() *Queue {
	q := &Queue{}
	q.buf = NewBuf()
	q.lock = &sync.RWMutex{}
	q.cond = sync.NewCond(q.lock.RLocker())
	q.dpidx = -1
	return q
}

// Close ----------------------------------------------------------------------
// After Close() called, all QueueCursor's ReadPacket will return io.EOF
// ----------------------------------------------------------------------------
func (que *Queue) Close() (err error) {
	que.lock.Lock()

	que.closed = true
	que.cond.Broadcast()

	que.lock.Unlock()
	return
}

// WritePacket ----------------------------------------------------------------
// Put packet into buffer, old packets will be discared
// ----------------------------------------------------------------------------
func (que *Queue) WritePacket(dp data.DataPoint) (err error) {
	que.lock.Lock()
	que.buf.Push(dp)
	que.cond.Broadcast()
	que.lock.Unlock()
	return
}


// Create cursor position at latest packet.
func (que *Queue) Latest() *QueueCursor {
	cursor := que.newCursor()
	cursor.init = func(buf *quBuf, dpidx int) BufPos {
		return buf.Tail
	}
	return cursor
}

// Oldest ---------------------------------------------------------------------
// Create cursor position at oldest buffered packet
// ----------------------------------------------------------------------------
func (que *Queue) Oldest() *QueueCursor {
	cursor := que.newCursor()
	cursor.init = func(buf *quBuf, dpidx int) BufPos {
		return buf.Head
	}
	return cursor
}

// ReadPacket -----------------------------------------------------------------
// will not consume packets in Queue, it's just a cursor
// ----------------------------------------------------------------------------
func (self *QueueCursor) ReadPacket() (dp data.DataPoint, err error) {
	self.que.cond.L.Lock()
	buf := self.que.buf
	if !self.gotpos {
		self.pos = self.init(buf, self.que.dpidx)
		self.gotpos = true
	}
	for {
		if self.pos.LT(buf.Head) {
			self.pos = buf.Head
		} else if self.pos.GT(buf.Tail) {
			self.pos = buf.Tail
		}
		if buf.IsValidPos(self.pos) {
			dp = buf.Get(self.pos)
			self.pos++
			break
		}
		if self.que.closed {
			err = io.EOF
			break
		}
		self.que.cond.Wait()
	}
	self.que.cond.L.Unlock()
	return
}
