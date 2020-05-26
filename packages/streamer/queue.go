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
	pos    BufIdx
	gotpos bool
	init   func(buf *quBuf, dpidx int) BufIdx
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
// Put packet into buffer, old packets will be discarded
// ----------------------------------------------------------------------------
func (que *Queue) WritePacket(dp data.DataPoint) (err error) {
	que.lock.Lock()
	que.buf.Push(dp)
	que.cond.Broadcast()
	que.lock.Unlock()
	return
}

// Latest ---------------------------------------------------------------------
// Create cursor position at latest (most recent) packet
// ----------------------------------------------------------------------------
func (que *Queue) Latest() *QueueCursor {
	cursor := que.newCursor()
	cursor.init = func(buf *quBuf, dpidx int) BufIdx {
		return buf.Tail
	}
	return cursor
}

// Oldest ---------------------------------------------------------------------
// Create cursor position at oldest buffered packet
// ----------------------------------------------------------------------------
func (que *Queue) Oldest() *QueueCursor {
	cursor := que.newCursor()
	cursor.init = func(buf *quBuf, dpidx int) BufIdx {
		return buf.Head
	}
	return cursor
}

// ReadPacket -----------------------------------------------------------------
// will not consume packets in Queue, it's just a cursor
// ----------------------------------------------------------------------------
func (qc *QueueCursor) ReadPacket() (dp data.DataPoint, err error) {
	qc.que.cond.L.Lock()
	buf := qc.que.buf
	if !qc.gotpos {
		qc.pos = qc.init(buf, qc.que.dpidx)
		qc.gotpos = true
	}
	for {
		if qc.pos.LT(buf.Head) {
			qc.pos = buf.Head
		} else if qc.pos.GT(buf.Tail) {
			qc.pos = buf.Tail
		}
		if buf.IsValidPos(qc.pos) {
			dp = buf.Get(qc.pos)
			qc.pos++
			break
		}
		if qc.que.closed {
			err = io.EOF
			break
		}
		qc.que.cond.Wait()
	}
	qc.que.cond.L.Unlock()
	return
}
