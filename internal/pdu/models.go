package pdu

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"smpp-server-golang/internal/helpers"
)

const (
	HeaderSize = 16
)

const (
	BindTransceiver   = 0x00000009
	EnquireLink       = 0x00000015
	UnbindTransceiver = 0x00000006

	BindTransceiverResponse   = 0x80000009
	EnquireLinkResponse       = 0x80000015
	UnbindTransceiverResponse = 0x80000006
)

type Pdu struct {
	Header     *Header
	Body       *Body
	HeaderByte []byte
	BodyByte   []byte
}

type VisualizationObject[T any] struct {
	Bytes       []byte
	Value       T
	IsBigEndian bool
	Name        string
}

func NewHeader(commandId uint32, commandStatus uint32, sequenceNumber uint32) *Header {

	commandIdBytes, err := helpers.AnyToBytes(commandId)
	if err != nil {
		return nil
	}

	commandStatusBytes, err := helpers.AnyToBytes(commandStatus)
	if err != nil {
		return nil
	}

	sequenceNumberBytes, err := helpers.AnyToBytes(sequenceNumber)
	if err != nil {
		return nil
	}
	newHeader := &Header{}
	newHeader.CommandId = helpers.NewVisualizationObject[uint32](commandIdBytes, helpers.WithBigEndian[uint32](true))
	newHeader.CommandStatus = helpers.NewVisualizationObject[uint32](commandStatusBytes, helpers.WithBigEndian[uint32](true))
	newHeader.SequenceNumber = helpers.NewVisualizationObject[uint32](sequenceNumberBytes, helpers.WithBigEndian[uint32](true))

	return newHeader
}

func (p *Pdu) ParseHeader(headerData []byte) {
	p.HeaderByte = headerData
	p.Header = &Header{}
	p.Header.CommandLength = helpers.NewVisualizationObject[uint32](headerData[:4], helpers.WithName[uint32]("Command Length"), helpers.WithBigEndian[uint32](true))
	p.Header.CommandId = helpers.NewVisualizationObject[uint32](headerData[4:8], helpers.WithName[uint32]("Command Id"), helpers.WithBigEndian[uint32](true))
	p.Header.CommandStatus = helpers.NewVisualizationObject[uint32](headerData[8:12], helpers.WithName[uint32]("Command Status"), helpers.WithBigEndian[uint32](true))
	p.Header.SequenceNumber = helpers.NewVisualizationObject[uint32](headerData[12:16], helpers.WithName[uint32]("Sequence Number"), helpers.WithBigEndian[uint32](true))

	println("Header Hex:", hex.EncodeToString(headerData))

	p.Header.CommandLength.Print()
	p.Header.CommandId.Print()
	p.Header.CommandStatus.Print()
	p.Header.SequenceNumber.Print()
}

func GetBytes(p Pdu) []byte {

	bodyLength := len(p.BodyByte)
	responseSize := HeaderSize + bodyLength
	println("Total size:", responseSize)

	response := make([]byte, responseSize)

	// Command Length (writes to buffer[0:4])
	binary.BigEndian.PutUint32(response[0:], uint32(responseSize))

	// Command ID (writes to buffer[4:8])
	binary.BigEndian.PutUint32(response[4:], p.Header.CommandId.GetValue())

	// Command Status (writes to buffer[8:12])
	binary.BigEndian.PutUint32(response[8:], p.Header.CommandStatus.GetValue())

	// Sequence Number (writes to buffer[12:16])
	binary.BigEndian.PutUint32(response[12:], p.Header.SequenceNumber.GetValue())

	if bodyLength > 0 {
		copy(response[HeaderSize:], p.BodyByte)
	}

	println("Writing TCP back:", hex.EncodeToString(response))

	return response
}

func SendPdu(conn net.Conn, response []byte) {
	_, err := conn.Write(response)
	if err != nil {
		println("Write error:", err)
	}

	println("Successfully responded back:", hex.EncodeToString(response))
}

type Header struct {
	CommandLength  *helpers.VisualizationObject[uint32]
	CommandStatus  *helpers.VisualizationObject[uint32]
	CommandId      *helpers.VisualizationObject[uint32]
	SequenceNumber *helpers.VisualizationObject[uint32]
}

type Body struct {
}
