// Code generated by github.com/skycoin/skyencoder. DO NOT EDIT.

package coin

import (
	"errors"
	"math"

	"github.com/ness-network/privateness/src/cipher"
	"github.com/ness-network/privateness/src/cipher/encoder"
)

// encodeSizeTransactionInputs computes the size of an encoded object of type transactionInputs
func encodeSizeTransactionInputs(obj *transactionInputs) uint64 {
	i0 := uint64(0)

	// obj.In
	i0 += 4
	{
		i1 := uint64(0)

		// x1
		i1 += 32

		i0 += uint64(len(obj.In)) * i1
	}

	return i0
}

// encodeTransactionInputs encodes an object of type transactionInputs to a buffer allocated to the exact size
// required to encode the object.
func encodeTransactionInputs(obj *transactionInputs) ([]byte, error) {
	n := encodeSizeTransactionInputs(obj)
	buf := make([]byte, n)

	if err := encodeTransactionInputsToBuffer(buf, obj); err != nil {
		return nil, err
	}

	return buf, nil
}

// encodeTransactionInputsToBuffer encodes an object of type transactionInputs to a []byte buffer.
// The buffer must be large enough to encode the object, otherwise an error is returned.
func encodeTransactionInputsToBuffer(buf []byte, obj *transactionInputs) error {
	if uint64(len(buf)) < encodeSizeTransactionInputs(obj) {
		return encoder.ErrBufferUnderflow
	}

	e := &encoder.Encoder{
		Buffer: buf[:],
	}

	// obj.In maxlen check
	if len(obj.In) > 65535 {
		return encoder.ErrMaxLenExceeded
	}

	// obj.In length check
	if uint64(len(obj.In)) > math.MaxUint32 {
		return errors.New("obj.In length exceeds math.MaxUint32")
	}

	// obj.In length
	e.Uint32(uint32(len(obj.In)))

	// obj.In
	for _, x := range obj.In {

		// x
		e.CopyBytes(x[:])

	}

	return nil
}

// decodeTransactionInputs decodes an object of type transactionInputs from a buffer.
// Returns the number of bytes used from the buffer to decode the object.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
func decodeTransactionInputs(buf []byte, obj *transactionInputs) (uint64, error) {
	d := &encoder.Decoder{
		Buffer: buf[:],
	}

	{
		// obj.In

		ul, err := d.Uint32()
		if err != nil {
			return 0, err
		}

		length := int(ul)
		if length < 0 || length > len(d.Buffer) {
			return 0, encoder.ErrBufferUnderflow
		}

		if length > 65535 {
			return 0, encoder.ErrMaxLenExceeded
		}

		if length != 0 {
			obj.In = make([]cipher.SHA256, length)

			for z1 := range obj.In {
				{
					// obj.In[z1]
					if len(d.Buffer) < len(obj.In[z1]) {
						return 0, encoder.ErrBufferUnderflow
					}
					copy(obj.In[z1][:], d.Buffer[:len(obj.In[z1])])
					d.Buffer = d.Buffer[len(obj.In[z1]):]
				}

			}
		}
	}

	return uint64(len(buf) - len(d.Buffer)), nil
}

// decodeTransactionInputsExact decodes an object of type transactionInputs from a buffer.
// If the buffer not long enough to decode the object, returns encoder.ErrBufferUnderflow.
// If the buffer is longer than required to decode the object, returns encoder.ErrRemainingBytes.
func decodeTransactionInputsExact(buf []byte, obj *transactionInputs) error {
	if n, err := decodeTransactionInputs(buf, obj); err != nil {
		return err
	} else if n != uint64(len(buf)) {
		return encoder.ErrRemainingBytes
	}

	return nil
}
