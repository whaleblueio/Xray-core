package dispatcher

import (
	"github.com/juju/ratelimit"
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/buf"
	"github.com/whaleblueio/Xray-core/features/stats"
)

type SizeStatWriter struct {
	Counter stats.Counter
	Writer  buf.Writer
	Bucket  *ratelimit.Bucket
}

func (w *SizeStatWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	w.Counter.Add(int64(mb.Len()))
	if w.Bucket != nil {
		w.Bucket.Take(int64(mb.Len()))
	}
	return w.Writer.WriteMultiBuffer(mb)
}

func (w *SizeStatWriter) Close() error {
	return common.Close(w.Writer)
}

func (w *SizeStatWriter) Interrupt() {
	common.Interrupt(w.Writer)
}
