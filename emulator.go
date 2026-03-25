package main

import (
	"sync"

	"github.com/loeredami/ungo"
)

type Worker struct {
	PID      uint64
	IP       uint64
	emulator *Emulator
	stack    *ungo.Stack[uint64]
	runnin   bool
}

type Emulator struct {
	workers []*Worker
	memory  []byte

	mu *sync.RWMutex
}

func NewEmulator() *Emulator {
	em := &Emulator{
		workers: make([]*Worker, 0),
		memory:  make([]byte, 0),

		mu: &sync.RWMutex{},
	}

	return em
}

func (emu *Emulator) InitProgram(data []byte) (memID ungo.Optional[uint64]) {
	memID = emu.Allocate(uint64(len(data)))

	memID.IfPresent(func(mID uint64) {
		success := emu.WriteAt(mID, data)
		memID = ungo.If(success, memID, ungo.None[uint64]())
	})

	return
}
