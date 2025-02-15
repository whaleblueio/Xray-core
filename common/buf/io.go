package buf

import (
	rateLimit "github.com/juju/ratelimit"
	"github.com/whaleblueio/Xray-core/common"
	"github.com/whaleblueio/Xray-core/common/log"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

// Reader extends io.Reader with MultiBuffer.
type Reader interface {
	// ReadMultiBuffer reads content from underlying reader, and put it into a MultiBuffer.
	ReadMultiBuffer() (MultiBuffer, error)
}

// ErrReadTimeout is an error that happens with IO timeout.
var ErrReadTimeout = newError("IO timeout")

// TimeoutReader is a reader that returns error if Read() operation takes longer than the given timeout.
type TimeoutReader interface {
	ReadMultiBufferTimeout(time.Duration) (MultiBuffer, error)
}

// Writer extends io.Writer with MultiBuffer.
type Writer interface {
	// WriteMultiBuffer writes a MultiBuffer into underlying writer.
	WriteMultiBuffer(MultiBuffer) error
}

// WriteAllBytes ensures all bytes are written into the given writer.
func WriteAllBytes(writer io.Writer, payload []byte) error {
	for len(payload) > 0 {
		n, err := writer.Write(payload)
		if err != nil {
			return err
		}
		payload = payload[n:]
	}
	return nil
}

func isPacketReader(reader io.Reader) bool {
	_, ok := reader.(net.PacketConn)
	return ok
}

func NewLimitReader(reader io.Reader, speed int64) Reader {
	if mr, ok := reader.(Reader); ok {
		newError("NewLimitReader() is multiple Reader sequenceId:", common.GetSequenceId()).WriteToLog()
		return mr
	}
	if speed > 0 {
		bucket := rateLimit.NewBucketWithQuantum(time.Second, speed, speed)
		limitReader := rateLimit.Reader(reader, bucket)
		newError("NewLimitReader() is speed limit Reader sequenceId:", common.GetSequenceId()).WriteToLog()

		if isPacketReader(reader) {
			return &PacketReader{
				Reader: limitReader,
			}
		}

		_, isFile := reader.(*os.File)
		if !isFile && useReadv {
			if sc, ok := reader.(syscall.Conn); ok {
				rawConn, err := sc.SyscallConn()
				if err != nil {
					newError("failed to get sysconn").Base(err).WriteToLog()
				} else {
					return NewReadVReader(limitReader, rawConn)
				}
			}
		}

		return &SingleReader{
			Reader: limitReader,
		}
	}

	if isPacketReader(reader) {
		return &PacketReader{
			Reader: reader,
		}
	}

	_, isFile := reader.(*os.File)
	if !isFile && useReadv {
		if sc, ok := reader.(syscall.Conn); ok {
			rawConn, err := sc.SyscallConn()
			if err != nil {
				newError("failed to get sysconn").Base(err).WriteToLog()
			} else {
				return NewReadVReader(reader, rawConn)
			}
		}
	}

	return &SingleReader{
		Reader: reader,
	}
}

// NewReader creates a new Reader.
// The Reader instance doesn't take the ownership of reader.
func NewReader(reader io.Reader) Reader {
	if mr, ok := reader.(Reader); ok {
		return mr
	}

	if isPacketReader(reader) {
		return &PacketReader{
			Reader: reader,
		}
	}

	_, isFile := reader.(*os.File)
	if !isFile && useReadv {
		if sc, ok := reader.(syscall.Conn); ok {
			rawConn, err := sc.SyscallConn()
			if err != nil {
				newError("failed to get sysconn").Base(err).WriteToLog()
			} else {
				return NewReadVReader(reader, rawConn)
			}
		}
	}

	return &SingleReader{
		Reader: reader,
	}
}

// NewPacketReader creates a new PacketReader based on the given reader.
func NewPacketReader(reader io.Reader) Reader {
	if mr, ok := reader.(Reader); ok {
		return mr
	}
	return &PacketReader{
		Reader: reader,
	}
}
func NewPacketReaderWithRateLimiter(reader io.Reader, speed int64) Reader {
	if mr, ok := reader.(Reader); ok {
		return mr
	}
	bucket := rateLimit.NewBucketWithQuantum(time.Second, speed, speed)
	limitReader := rateLimit.Reader(reader, bucket)
	return &PacketReader{
		Reader: limitReader,
	}
}

func isPacketWriter(writer io.Writer) bool {
	if _, ok := writer.(net.PacketConn); ok {
		return true
	}

	// If the writer doesn't implement syscall.Conn, it is probably not a TCP connection.
	if _, ok := writer.(syscall.Conn); !ok {
		return true
	}
	return false
}

// NewWriter creates a new Writer.
func NewWriter(writer io.Writer) Writer {
	if mw, ok := writer.(Writer); ok {
		return mw
	}

	if isPacketWriter(writer) {
		return &SequentialWriter{
			Writer: writer,
		}
	}

	return &BufferToBytesWriter{
		Writer: writer,
	}
}

// NewWriter creates a new Writer.
func NewWriterWithRateLimiter(writer io.Writer, speed int64) Writer {
	if mw, ok := writer.(Writer); ok {
		log.Record(&log.GeneralMessage{
			Content: "NewWriterWithRateLimiter() is multiple writer",
		})
		return mw
	}
	if speed > 0 {
		bucket := rateLimit.NewBucketWithQuantum(time.Second, speed, speed)
		limitWriter := rateLimit.Writer(writer, bucket)
		if isPacketWriter(writer) {
			return &SequentialWriter{
				Writer: limitWriter,
			}
		}

		return &BufferToBytesWriter{
			Writer: limitWriter,
		}
	}
	if isPacketWriter(writer) {
		return &SequentialWriter{
			Writer: writer,
		}
	}

	return &BufferToBytesWriter{
		Writer: writer,
	}
}
