package pdu

import (
	"encoding/hex"
	"smpp-server-golang/internal/helpers"
)

const (
	HeaderSize = 16
)

const (
	BindReceiver            = 0x80000000
	BindTransceiver         = 0x00000009
	EnquireLink             = 0x00000015
	BindReceiverResponse    = 0x80000001
	BindTransceiverResponse = 0x80000009
	EnquireLinkResponse     = 0x80000015
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
	p.Header.CommandLength = helpers.NewVisualizationObject[uint32](headerData[:4], "Command Length", helpers.WithBigEndian[uint32](true))
	p.Header.CommandId = helpers.NewVisualizationObject[uint32](headerData[4:8], "Command Id", helpers.WithBigEndian[uint32](true))
	p.Header.CommandStatus = helpers.NewVisualizationObject[uint32](headerData[8:12], "Command Status", helpers.WithBigEndian[uint32](true))
	p.Header.SequenceNumber = helpers.NewVisualizationObject[uint32](headerData[12:16], "Sequence Number", helpers.WithBigEndian[uint32](true))

	println("Header Hex:", hex.EncodeToString(headerData))

	p.Header.CommandLength.Print()
	p.Header.CommandId.Print()
	p.Header.CommandStatus.Print()
	p.Header.SequenceNumber.Print()
}

type Header struct {
	CommandLength  *helpers.VisualizationObject[uint32]
	CommandStatus  *helpers.VisualizationObject[uint32]
	CommandId      *helpers.VisualizationObject[uint32]
	SequenceNumber *helpers.VisualizationObject[uint32]
}

type Body struct {
}
