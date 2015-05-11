package main

import (
	"github.com/edmccard/avr-sim/core"
	"github.com/edmccard/ihex"
	"io"
)

type MemRead func(core.Addr) byte
type MemWrite func(core.Addr, byte)

type DemoMem struct {
	prog     []uint16
	data     []byte
	inports  []MemRead
	outports []MemWrite
	dmask    core.Addr
	pmask    core.Addr
}

func NewDemoMem(cpu *core.Cpu, data io.Reader) *DemoMem {
	mem := &DemoMem{
		prog:     make([]uint16, 0x4000),
		data:     make([]byte, 0x0500),
		inports:  make([]MemRead, 0x60),
		outports: make([]MemWrite, 0x60),
		dmask:    0xffff,
		pmask:    0xfffff,
	}
	readreg := func(addr core.Addr) byte {
		return cpu.GetReg(int(addr))
	}
	writereg := func(addr core.Addr, val byte) {
		cpu.SetReg(int(addr), val)
	}
	readspl := func(addr core.Addr) byte {
		return cpu.GetSPL()
	}
	writespl := func(addr core.Addr, val byte) {
		cpu.SetSPL(val)
	}
	readsph := func(addr core.Addr) byte {
		return cpu.GetSPH()
	}
	writesph := func(addr core.Addr, val byte) {
		cpu.SetSPH(val)
	}
	readsreg := func(addr core.Addr) byte {
		return cpu.ByteFromSreg()
	}
	writesreg := func(addr core.Addr, val byte) {
		cpu.SregFromByte(val)
	}
	readmem := func(addr core.Addr) byte {
		return mem.data[addr]
	}
	writemem := func(addr core.Addr, val byte) {
		mem.data[addr] = val
	}
	for i := 0; i < 32; i++ {
		mem.inports[i] = readreg
		mem.outports[i] = writereg
	}
	for i := 32; i < 96; i++ {
		mem.inports[i] = readmem
		mem.outports[i] = writemem
	}
	mem.inports[0x5d] = readspl
	mem.outports[0x5d] = writespl
	mem.inports[0x5e] = readsph
	mem.outports[0x5e] = writesph
	mem.inports[0x5f] = readsreg
	mem.outports[0x5f] = writesreg

	mem.loadHex(data)

	return mem
}

func (mem *DemoMem) loadHex(data io.Reader) {
	parser := ihex.NewParser(data)
	for parser.Parse() {
		mem.loadRecord(parser.Data())
	}
	if parser.Err() != nil {
		panic("bad hex data")
	}
}

func (mem *DemoMem) loadRecord(rec ihex.Record) {
	addr := rec.Address >> 1
	for i := 0; i < len(rec.Bytes); i += 2 {
		val := uint16(rec.Bytes[i]) | (uint16(rec.Bytes[i+1]) << 8)
		mem.prog[addr] = val
		addr++
	}
}

func (mem *DemoMem) LoadProgram(addr core.Addr) byte {
	shift := (uint(addr) & 0x1) * 8
	addr = (addr >> 1) & mem.pmask
	return byte(mem.prog[addr] >> shift)
}

func (mem *DemoMem) ReadProgram(addr core.Addr) uint16 {
	addr &= mem.pmask
	return mem.prog[addr]
}

func (mem *DemoMem) ReadData(addr core.Addr) byte {
	addr &= mem.dmask
	if addr < 0x60 {
		return mem.inports[addr](addr)
	}
	return mem.data[addr]
}

func (mem *DemoMem) WriteData(addr core.Addr, val byte) {
	addr &= mem.dmask
	if addr < 0x60 {
		mem.outports[addr](addr, val)
		return
	}
	mem.data[addr] = val
}
