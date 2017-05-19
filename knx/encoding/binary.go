package encoding

import (
	"encoding/binary"
	"errors"
	"io"
)

func writeUint8(w io.Writer, buffer []byte, value uint8) (int, error) {
	buffer[0] = value
	return w.Write(buffer[:1])
}

func writeUint16(w io.Writer, buffer []byte, value uint16) (int, error) {
	binary.BigEndian.PutUint16(buffer[:2], value)
	return w.Write(buffer[:2])
}

func writeUint32(w io.Writer, buffer []byte, value uint32) (int, error) {
	binary.BigEndian.PutUint32(buffer[:4], value)
	return w.Write(buffer[:4])
}

func writeUint64(w io.Writer, buffer []byte, value uint64) (int, error) {
	binary.BigEndian.PutUint64(buffer[:8], value)
	return w.Write(buffer[:8])
}

func writeUint16Slice(w io.Writer, values []uint16) (int, error) {
	buffer := make([]byte, 2 * len(values))

	for i, value := range values {
		binary.BigEndian.PutUint16(buffer[i * 2:], value)
	}

	return w.Write(buffer)
}

func writeUint32Slice(w io.Writer, values []uint32) (int, error) {
	buffer := make([]byte, 4 * len(values))

	for i, value := range values {
		binary.BigEndian.PutUint32(buffer[i * 4:], value)
	}

	return w.Write(buffer)
}

func writeUint64Slice(w io.Writer, values []uint64) (int, error) {
	buffer := make([]byte, 8 * len(values))

	for i, value := range values {
		binary.BigEndian.PutUint64(buffer[i * 8:], value)
	}

	return w.Write(buffer)
}

//
var (
	ErrUnknownDataType = errors.New("Unknown data type can't be written")
)

func writeItem(w io.Writer, item interface{}) (n int, err error) {
	buffer := make([]byte, 8)

	switch item := item.(type) {
	case uint8:  return writeUint8(w, buffer, item)
	case uint16: return writeUint16(w, buffer, item)
	case uint32: return writeUint32(w, buffer, item)
	case uint64: return writeUint64(w, buffer, item)
	case int8:   return writeUint8(w, buffer, uint8(item))
	case int16:  return writeUint16(w, buffer, uint16(item))
	case int32:  return writeUint32(w, buffer, uint32(item))
	case int64:  return writeUint64(w, buffer, uint64(item))

	case *uint8:  return writeUint8(w, buffer, *item)
	case *uint16: return writeUint16(w, buffer, *item)
	case *uint32: return writeUint32(w, buffer, *item)
	case *uint64: return writeUint64(w, buffer, *item)
	case *int8:   return writeUint8(w, buffer, uint8(*item))
	case *int16:  return writeUint16(w, buffer, uint16(*item))
	case *int32:  return writeUint32(w, buffer, uint32(*item))
	case *int64:  return writeUint64(w, buffer, uint64(*item))

	case []uint8:  return w.Write(item)
	case []uint16: return writeUint16Slice(w, item)
	case []uint32: return writeUint32Slice(w, item)
	case []uint64: return writeUint64Slice(w, item)

	case io.WriterTo:
		len, err := item.WriteTo(w)
		return int(len), err

	default:
		return 0, ErrUnknownDataType
	}
}

// WriteSequence writes items to the Writer.
func WriteSequence(w io.Writer, items ...interface{}) (n int, err error) {
	for _, item := range items {
		len, err := writeItem(w, item)
		n += len
		if err != nil {
			return n, err
		}
	}

	return n, err
}

// ReadSequence reads items from the Reader.
func ReadSequence(r io.Reader, items ...interface{}) error {
	for _, item := range items {
		err := binary.Read(r, binary.BigEndian, item)
		if err != nil { return err }
	}

	return nil
}

// UInt16 extracts a uint16 from the given data.
func UInt16(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}