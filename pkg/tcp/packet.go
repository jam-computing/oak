package tcp

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net"
)

var PROTOCOL_VERSION byte = 2

type FullPacket struct {
	Meta *PacketMeta
	Data *PacketData
	Id   uint16
}

func NewFullPacket(meta PacketMeta, data *string, number *uint16) *FullPacket {
	var r int
	if number == nil {
		r = rand.IntN(5000)
	}
	p := FullPacket{
		Meta: &meta,
		Data: nil,
		Id:   uint16(r),
	}
	if data != nil {
		p.Data = &PacketData{
			Version: PROTOCOL_VERSION,
			Id:      uint16(r),
			Data:    *data,
		}
	}
	p.Meta.Id = p.Id
	if p.Data != nil {
		p.Data.Id = p.Id
        p.Meta.Len = uint16(len(p.Data.Data))
	}
	return &p
}

type PacketMeta struct {
	Version byte
	Id      uint16
	Command byte
	Status  uint16
	Len     uint16
}

type PacketData struct {
	Version byte
	Id      uint16
	Data    string
}

func NewMetaPacket() PacketMeta {
	packet := PacketMeta{}

	packet.Version = PROTOCOL_VERSION
	packet.Id = 0
	packet.Command = 2
	packet.Status = 200
	packet.Len = 0

	return packet
}

func BadPacket() PacketMeta {
	packet := PacketMeta{}

	packet.Version = PROTOCOL_VERSION
	packet.Command = 0
	packet.Status = 500
	packet.Id = 0
	packet.Len = 0

	return packet
}

func MakeMeta(bytes []byte) PacketMeta {
	packet := NewMetaPacket()

	packet.Version = bytes[0]
	packet.Id = binary.LittleEndian.Uint16([]byte{bytes[1], bytes[2]})
	packet.Command = bytes[3]
	packet.Status = binary.LittleEndian.Uint16([]byte{bytes[4], bytes[5]})
	packet.Len = binary.LittleEndian.Uint16([]byte{bytes[6], bytes[7]})

	return packet
}

func MakeData(bytes []byte) *PacketData {
	packet := PacketData{}
	packet.Version = bytes[0]

	if packet.Version != PROTOCOL_VERSION {
		return nil
	}

	packet.Id = binary.LittleEndian.Uint16([]byte{bytes[1], bytes[2]})
	packet.Data = string(bytes[3:])

	return &packet
}

func (packet *PacketMeta) ToBytes() []byte {
	status := make([]byte, 2)
	len := make([]byte, 2)
	id := make([]byte, 2)
	bytes := make([]byte, 0, 8)

	binary.LittleEndian.PutUint16(status, packet.Status)
	binary.LittleEndian.PutUint16(id, packet.Id)
	binary.LittleEndian.PutUint16(len, packet.Len)

	bytes = append(bytes, packet.Version)
	bytes = append(bytes, id...)
	bytes = append(bytes, packet.Command)
	bytes = append(bytes, status...)
	bytes = append(bytes, len...)

	return bytes
}

func (packet *PacketData) DataToBytes(meta *PacketMeta) []byte {
	bytes := make([]byte, 0, meta.Len+3)
	id := make([]byte, 2)

	binary.LittleEndian.PutUint16(id, packet.Id)

	bytes = append(bytes, packet.Version)
	bytes = append(bytes, id...)
	bytes = append(bytes, []byte(packet.Data)...)

	return bytes
}

func (packet FullPacket) SendRecv() *FullPacket {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println(err)
		return nil
	}

	_, err = conn.Write(packet.Meta.ToBytes())

	if err != nil {
		return nil
	}

	if packet.Meta.Len > 0 {
		_, err = conn.Write(packet.Data.DataToBytes(packet.Meta))
		if err != nil {
			return nil
		}
	}

	// Read metadata packet
	metaBuf := make([]byte, 8) // Assuming metadata is always 8 bytes
	_, err = conn.Read(metaBuf)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	metaPacket := MakeMeta(metaBuf)


	// Read data packet based on the length from metadata
	if metaPacket.Len > 0 {
		dataBuf := make([]byte, metaPacket.Len+3)
		_, err = conn.Read(dataBuf)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		dataPacket := MakeData(dataBuf)
		conn.Close()
		return NewFullPacket(metaPacket, &dataPacket.Data, &metaPacket.Id)
	}

	return NewFullPacket(metaPacket, nil, nil)

}
