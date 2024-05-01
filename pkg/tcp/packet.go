package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
)

var PROTOCOL_VERSION byte = 1

type Packet struct {
	Version byte
	Command byte
	Status  uint16
	Id      uint16
	Len     uint16
	Data    string
}

func NewPacket() Packet {
	packet := Packet{}

	packet.Version = PROTOCOL_VERSION
	packet.Command = 1
	packet.Status = 200
	packet.Id = 0
	packet.Len = 5
	packet.Data = "Hello"

	return packet
}

func BadPacket() Packet {
	packet := Packet{}

	packet.Version = PROTOCOL_VERSION
	packet.Command = 0
	packet.Status = 500
	packet.Id = 0
    packet.Len = 0
	packet.Data = ""

	return packet
}

func Make(bytes []byte) Packet {
	packet := NewPacket()

	packet.Version = bytes[0]
	packet.Command = bytes[1]

	packet.Status = binary.LittleEndian.Uint16([]byte{bytes[2], bytes[3]})
	packet.Id = binary.LittleEndian.Uint16([]byte{bytes[4], bytes[5]})
	packet.Len = binary.LittleEndian.Uint16([]byte{bytes[6], bytes[7]}) // Length
	packet.Data = string(bytes[8:])

	return packet
}

func (packet *Packet) ToBytes() []byte {

	// This feels bad
	bytes := make([]byte, 0, 1024)
	bytes = append(bytes, packet.Version)
	bytes = append(bytes, packet.Command)

	status := make([]byte, 2)
	id := make([]byte, 2)

	binary.LittleEndian.PutUint16(status, packet.Status)
	binary.LittleEndian.PutUint16(id, packet.Id)

	bytes = append(bytes, status...)
	bytes = append(bytes, id...)

	if len(packet.Data) > 0 {
		data_len := uint16(len(packet.Data))
		len := make([]byte, 2)
		binary.LittleEndian.PutUint16(len, data_len)

		bytes = append(bytes, len...)
		bytes = append(bytes, packet.Data...)
	} else {
		len := make([]byte, 2)
		binary.LittleEndian.PutUint16(len, 0)
		bytes = append(bytes, len...)
	}

	return bytes
}

func (packet *Packet) SendRecv() Packet {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println(err)
		return BadPacket()
	}

	_, err = conn.Write(packet.ToBytes())

	if err != nil {
		fmt.Println(err)
		return BadPacket()
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)

	if err != nil {
		fmt.Println(err)
		return BadPacket()
	}

	conn.Close()
	return Make(buf)
}

