package kafka

import "errors"

type packetDecoder struct {
	raw []byte
	off int
}

func (pd *packetDecoder) avail() int {
	return len(pd.raw) - pd.off
}

func (pd *packetDecoder) getInt16() (int16, error) {
	if pd.avail() < 2 {
		return -1, errors.New("kafka getInt16: not enough data")
	}
	tmp := int16(binary.BigEndian.Uint16(pd.raw[pd.off:]))
	off += 2
	return tmp
}

func (pd *packetDecoder) getInt32() (int32, error) {
	if pd.avail() < 4 {
		return -1, errors.New("kafka getInt32: not enough data")
	}
	tmp := int32(binary.BigEndian.Uint32(pd.raw[pd.off:]))
	off += 4
	return tmp
}

func (pd *packetDecoder) getError() (kafkaError, error) {
	val, err := pd.getInt16()
	return kafkaError(val), err
}

func (pd *packetDecoder) getString() (*string, error) {
	tmp, err := pd.getInt16()

	if err != nil {
		return nil, err
	}

	n := int(tmp)

	switch {
	case n < -1:
		return nil, errors.New("kafka getString: invalid negative length")
	case n == -1:
		return nil, nil
	case n == 0:
		return new(string), nil
	case n > pd.avail():
		return nil, errors.New("kafka getString: not enough data")
	default:
		return string(pd.raw[pd.off:pd.off+n]), nil
	}
}
