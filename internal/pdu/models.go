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

func (p *Pdu) Reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (p *Pdu) GetBytes(conn net.Conn, sequenceNumber uint32) {
	//0000001180000009000000000000000100
	responseSize := HeaderSize + len(p.BodyByte)
	println("Total size:", responseSize)

	response := make([]byte, responseSize)
	response = append(response, byte(responseSize))
	p.Reverse(response)
	println("CommandLength:", hex.EncodeToString(response))
	binary.BigEndian.PutUint32(response, BindReceiverResponse)
	println("CommandId:", hex.EncodeToString(response))
	response = append(response, 0)
	println("CommandStatus:", hex.EncodeToString(response))
	binary.LittleEndian.PutUint32(response[4:], sequenceNumber)
	println("SequenceNumber:", hex.EncodeToString(response))
	//wholeBody := append(p.BodyByte, "\000"...)
	//response = append(response, wholeBody...)

	println("Writing TCP back:", hex.EncodeToString(response))

	_, err := conn.Write(response)
	if err != nil {
		println("Write error:", err)
	}

}

type Header struct {
	CommandLength  *helpers.VisualizationObject[uint32]
	CommandStatus  *helpers.VisualizationObject[uint32]
	CommandId      *helpers.VisualizationObject[uint32]
	SequenceNumber *helpers.VisualizationObject[uint32]
}

type Body struct {
}
