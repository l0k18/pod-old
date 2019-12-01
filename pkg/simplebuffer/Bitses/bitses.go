package Bitses

import (
	"encoding/binary"
	"sync"
)

type Bitses struct {
	sync.Mutex
	Length  byte
	Byteses map[int32][]byte
}

func NewBitses() *Bitses {
	return &Bitses{Byteses: make(map[int32][]byte)}
}

func (b *Bitses) DecodeOne(by []byte) *Bitses {
	b.Decode(by)
	return b
}

func (b *Bitses) Decode(by []byte) (out []byte) {
	b.Lock()
	defer b.Unlock()
	//log.SPEW(by)
	if len(by) >= 7 {
		nB := by[0]
		if len(by) >= int(nB)*8 {
			for i := 0; i < int(nB); i++ {
				algoVer := int32(binary.BigEndian.Uint32(by[1+i*8 : 1+i*8+4]))
				//log.DEBUG("algoVer", algoVer, by[1+i*8+4:1+i*8+8], b.Byteses)
				b.Byteses[algoVer] = by[1+i*8+4 : 1+i*8+8]
			}
		}
		bL := int(nB)*8 + 1
		if len(by) > bL {
			out = by[bL:]
		}
	}
	//log.SPEW(b.Byteses)
	return
}

func (b *Bitses) Encode() (out []byte) {
	b.Lock()
	defer b.Unlock()
	out = []byte{b.Length}
	for algoVer := range b.Byteses {
		by := make([]byte, 4)
		binary.BigEndian.PutUint32(by, uint32(algoVer))
		out = append(out, append(by, b.Byteses[algoVer]...)...)
	}
	//log.SPEW(out)
	return
}

func (b *Bitses) Get() (out map[int32]uint32) {
	b.Lock()
	defer b.Unlock()
	out = make(map[int32]uint32)
	for algoVer := range b.Byteses {
		oB := binary.BigEndian.Uint32(b.Byteses[algoVer])
		out[algoVer] = oB
	}
	//log.SPEW(out)
	return
}

func (b *Bitses) Put(in map[int32]uint32) *Bitses {
	b.Lock()
	defer b.Unlock()
	b.Length = byte(len(in))
	b.Byteses = make(map[int32][]byte, b.Length)
	for algoVer := range in {
		bits := make([]byte, 4)
		binary.BigEndian.PutUint32(bits, in[algoVer])
		b.Byteses[algoVer] = bits
	}
	//log.SPEW(b.Byteses)
	return b
}
