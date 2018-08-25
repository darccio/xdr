// Copyright (C) 2014 Jakob Borg. All rights reserved.
// Copyright (C) 2018 Dario Castañé. All rights reserved. Use of this source code
// is governed by an MIT-style license that can be found in the LICENSE file.

package xdr

import "io"

// Unmarshaller is a thin wrapper around a byte buffer. The Unmarshal... methods
// don't individually return an error - the intention is that multiple fields are
// unmarshalled in rapid succession, followed by a check of the Error field on
// the Unmarshaller.
type Unmarshaller struct {
	Error error
	Data  []byte
}

// UnmarshalRaw returns a byte slice of length l from the buffer,
// without a size prefix or padding. This is suitable for retrieving
// data already in XDR format.
func (u *Unmarshaller) UnmarshalRaw(l int) []byte {
	if u.Error != nil {
		return nil
	}
	if len(u.Data) < l {
		u.Error = io.ErrUnexpectedEOF
		return nil
	}

	v := u.Data[:l]
	u.Data = u.Data[l:]

	return v
}

// UnmarshalString returns a string from the buffer.
func (u *Unmarshaller) UnmarshalString() string {
	return u.UnmarshalStringMax(0)
}

// UnmarshalStringMax returns a string up to a max length from the buffer.
func (u *Unmarshaller) UnmarshalStringMax(max int) string {
	buf := u.UnmarshalBytesMax(max)
	if len(buf) == 0 || u.Error != nil {
		return ""
	}

	return string(buf)
}

// UnmarshalBytes returns a byte slice from the buffer.
func (u *Unmarshaller) UnmarshalBytes() []byte {
	return u.UnmarshalBytesMax(0)
}

// UnmarshalBytesMax returns a byte slice up to a max length from the buffer.
func (u *Unmarshaller) UnmarshalBytesMax(max int) []byte {
	if u.Error != nil {
		return nil
	}
	if len(u.Data) < 4 {
		u.Error = io.ErrUnexpectedEOF
		return nil
	}

	l := int(u.Data[3]) | int(u.Data[2])<<8 | int(u.Data[1])<<16 | int(u.Data[0])<<24
	if l == 0 {
		u.Data = u.Data[4:]
		return nil
	}
	if l < 0 || max > 0 && l > max {
		// l may be negative on 32 bit builds
		u.Error = ElementSizeExceeded("bytes field", l, max)
		return nil
	}
	if len(u.Data) < l+4 {
		u.Error = io.ErrUnexpectedEOF
		return nil
	}

	v := u.Data[4 : 4+l]
	u.Data = u.Data[4+l+Padding(l):]

	return v
}

// UnmarshalBool returns a bool from the buffer.
func (u *Unmarshaller) UnmarshalBool() bool {
	return u.UnmarshalUint8() != 0
}

// UnmarshalUint8 returns a uint8 from the buffer.
func (u *Unmarshaller) UnmarshalUint8() uint8 {
	if u.Error != nil {
		return 0
	}
	if len(u.Data) < 4 {
		u.Error = io.ErrUnexpectedEOF
		return 0
	}

	v := uint8(u.Data[3])
	u.Data = u.Data[4:]

	return v
}

// UnmarshalUint16 returns a uint16 from the buffer.
func (u *Unmarshaller) UnmarshalUint16() uint16 {
	if u.Error != nil {
		return 0
	}
	if len(u.Data) < 4 {
		u.Error = io.ErrUnexpectedEOF
		return 0
	}

	v := uint16(u.Data[3]) | uint16(u.Data[2])<<8
	u.Data = u.Data[4:]

	return v
}

// UnmarshalUint32 returns a uint32 from the buffer.
func (u *Unmarshaller) UnmarshalUint32() uint32 {
	if u.Error != nil {
		return 0
	}
	if len(u.Data) < 4 {
		u.Error = io.ErrUnexpectedEOF
		return 0
	}

	v := uint32(u.Data[3]) | uint32(u.Data[2])<<8 | uint32(u.Data[1])<<16 | uint32(u.Data[0])<<24
	u.Data = u.Data[4:]

	return v
}

// UnmarshalUint64 returns a uint64 from the buffer.
func (u *Unmarshaller) UnmarshalUint64() uint64 {
	if u.Error != nil {
		return 0
	}
	if len(u.Data) < 8 {
		u.Error = io.ErrUnexpectedEOF
		return 0
	}

	v := uint64(u.Data[7]) | uint64(u.Data[6])<<8 | uint64(u.Data[5])<<16 | uint64(u.Data[4])<<24 |
		uint64(u.Data[3])<<32 | uint64(u.Data[2])<<40 | uint64(u.Data[1])<<48 | uint64(u.Data[0])<<56
	u.Data = u.Data[8:]

	return v
}
