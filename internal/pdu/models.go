package pdu

import (
	"encoding/binary"
	"encoding/hex"
	"smpp-server-golang/internal/helpers"
)

const (
	HeaderSize = 16
)

type Pdu struct {
	Header     *Header
	Body       *Body
	HeaderByte []byte
	BodyByte   []byte
}

func (p *Pdu) ParseHeader(headerData []byte) {
	p.HeaderByte = headerData
	p.Header = &Header{}
	p.Header.commandLength = binary.BigEndian.Uint32(headerData[:4])

	println("Header Hex:", hex.EncodeToString(headerData))
	println("Command Length:", p.Header.commandLength, "and Hex:", hex.EncodeToString(headerData[:4]))
}

type Header struct {
	commandLength  VisualizationObject[uint32]
	commandStatus  uint32
	commandId      uint32
	sequenceNumber uint32
}

type Body struct {
}

type VisualizationObject[T any] struct {
	Bytes       []byte
	Value       T
	IsBigEndian bool
}

func (v VisualizationObject[T]) Print() {
	println("Value", helpers.If[v.Value](v.IsBigEndian, binary.BigEndian.Uint32(v.Value), string(v.Value)))
}

type Print interface {
	Print()
}
