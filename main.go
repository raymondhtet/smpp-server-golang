package main

import (
	"net"
	"smpp-server-golang/internal/pdu"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	listener, err := net.Listen("tcp", ":2775")
	if err != nil {
		println("TCP Error listening at port 2775", err.Error())
		return
	}
	println("SMPP Server Listening at port 2775")

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			println("TCP Error closing listener at port 2775", err.Error())
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accepting connection", err.Error())
		}
		println("Accepted new connection")
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	println("New connection from " + conn.RemoteAddr().String())
	req := pdu.Pdu{}

	headerByte := make([]byte, pdu.HeaderSize)
	_, err := conn.Read(headerByte)
	if err != nil {
		println("Error reading total packet size", err.Error())
		return
	}
	req.ParseHeader(headerByte)

	switch req.Header.CommandId.Value {
	case pdu.BindTransceiver:
		response := pdu.Pdu{}
		response.GetBytes(conn, req.Header.SequenceNumber.Value)
		println("Bind Transceiver")
	case pdu.EnquireLink:
		println("EnquireLink")
	default:
		println("Unsupported CommandId")
	}

	//totalPacketLength := binary.BigEndian.Uint32(headerByte)
	//
	//println("Total Packet received:", hex.EncodeToString(headerByte), " and value:", totalPacketLength)
	//
	//pduSize := totalPacketLength - constants.HeaderSize
	//
	//println("PDU Size: ", pduSize)
	//
	//pduBytes := make([]byte, pduSize)
	//reader := bufio.NewReader(conn)
	//_, err = io.ReadFull(reader, pduBytes)
	//
	//if err != nil {
	//	println("Error reading full size", err.Error())
	//	return
	//}
	//
	//println("PDU received: ", hex.EncodeToString(pduBytes))
	//
	//wholePdu := append(headerByte, pduBytes...)
	//
	//println("Whole PDU received: ", hex.EncodeToString(wholePdu))
}
