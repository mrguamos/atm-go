package main

import (
	"fmt"
	"strings"

	"github.com/moov-io/iso8583/prefix"
)

type bcdVarPrefixer struct {
	Digits int
}

type bcdFixedPrefixer struct {
}

var BCDPrefixer = prefix.Prefixers{
	Fixed: &bcdFixedPrefixer{},
	L:     &bcdVarPrefixer{2},
	LL:    &bcdVarPrefixer{2},
	LLL:   &bcdVarPrefixer{3},
	LLLL:  &bcdVarPrefixer{4},
}

type bcdEncoder struct{}

var (
	// _          encoding.Encoder = (*bcdEncoder)(nil)
	BCDEncoder = &bcdEncoder{}
)

func (p *bcdVarPrefixer) EncodeLength(maxLen, dataLen int) ([]byte, error) {
	if dataLen > maxLen {
		return nil, fmt.Errorf("field length: %d is larger than maximum: %d", dataLen, maxLen)
	}
	var alloc int
	if p.Digits <= 2 {
		alloc = 1
	} else {
		alloc = 2
	}
	length := make([]byte, alloc)
	if alloc == 1 {
		length[0] = uint8(dataLen)
	} else {
		length[0] = byte(dataLen >> 8)
		length[1] = uint8(dataLen)
	}

	return length, nil
}

func (p *bcdVarPrefixer) DecodeLength(maxLen int, data []byte) (int, int, error) {
	var alloc int
	if p.Digits <= 2 {
		alloc = 1
	} else {
		alloc = 2
	}
	var dataLen int
	if alloc == 1 {
		dataLen = int(data[0])
	} else {
		dataLen = int(data[0])<<8 + int(data[1])
	}

	return dataLen, alloc, nil
}

func (p *bcdVarPrefixer) Inspect() string {
	return fmt.Sprintf("BCD.%s", strings.Repeat("L", p.Digits))
}

func (p *bcdFixedPrefixer) EncodeLength(fixLen, dataLen int) ([]byte, error) {
	if dataLen > fixLen {
		return nil, fmt.Errorf("field length: %d should be fixed: %d", dataLen, fixLen)
	}

	return []byte{}, nil
}

// Returns number of characters that should be decoded
func (p *bcdFixedPrefixer) DecodeLength(fixLen int, data []byte) (int, int, error) {
	return fixLen, 0, nil
}

func (p *bcdFixedPrefixer) Inspect() string {
	return "BCD.Fixed"
}
