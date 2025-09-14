package pdu

type Parser interface {
	ParseHeader(headerData []byte)
}
