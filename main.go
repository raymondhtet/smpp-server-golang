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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			println("Error closing connection", err.Error())
		}
	}(conn)

	println("New connection from " + conn.RemoteAddr().String())

	for {
		req := pdu.Pdu{}

		// Read header
		headerByte := make([]byte, pdu.HeaderSize)
		_, err := conn.Read(headerByte)
		if err != nil {
			println("Error reading total packet size", err.Error())
			return
		}
		req.ParseHeader(headerByte)

		/*
			// Read body if present
			commandLength := req.Header.CommandLength.GetValue()
			bodyLength := int(commandLength) - pdu.HeaderSize

			if bodyLength > 0 {
				bodyByte := make([]byte, bodyLength)
				_, err := conn.Read(bodyByte)
				if err != nil {
					println("Error reading body:", err.Error())
					return
				}
				req.BodyByte = bodyByte
				println("Read body of", bodyLength, "bytes")
			}
		*/

		response := pdu.Pdu{}
		println("Received Command ID:", req.Header.CommandId.GetValue())

		switch req.Header.CommandId.GetValue() {
		case pdu.BindTransceiver:
			println("Receiving Bind Transceiver Command")
			response.Header = pdu.NewHeader(pdu.BindTransceiverResponse, 0, req.Header.SequenceNumber.GetValue())
			response.BodyByte = []byte("\000")

		case pdu.EnquireLink:
			println("Receiving Enquire Link Command")
			// 0000001180000015000000000000000500 -> golang
			// 00000010800000150000000000000002 ->
			response.Header = pdu.NewHeader(pdu.EnquireLinkResponse, 0, req.Header.SequenceNumber.GetValue())

		case pdu.UnbindTransceiver:
			println("Receiving Unbind Transceiver Command")
			response.Header = pdu.NewHeader(pdu.UnbindTransceiverResponse, 0, req.Header.SequenceNumber.GetValue())

			responseBytes := pdu.GetBytes(response)
			pdu.SendPdu(conn, responseBytes)
			return

		default:
			println("Unsupported CommandId")
			continue // Skip sending response for unsupported commands
		}

		// Only send response if header was created
		responseBytes := pdu.GetBytes(response)
		pdu.SendPdu(conn, responseBytes)

	}
}
