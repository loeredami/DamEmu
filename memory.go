package main

import (
	"encoding/binary"

	"github.com/loeredami/ungo"
)

const MEM_HEADER_SIZE = uint64(16)

func (emu *Emulator) Allocate(size uint64) (mID ungo.Optional[uint64]) {
	if emu.mu.TryLock() {
		defer emu.mu.Unlock()
	}

	mID = ungo.None[uint64]()

	if size == 0 {
		return
	}

	highestID := uint64(0)
	scanLoc := uint64(0)
	for scanLoc+MEM_HEADER_SIZE <= uint64(len(emu.memory)) {
		header := binary.LittleEndian.Uint64(emu.memory[scanLoc : scanLoc+8])
		if header != 0 {
			if header > highestID {
				highestID = header
			}
			blockSize := binary.LittleEndian.Uint64(emu.memory[scanLoc+8 : scanLoc+16])
			scanLoc += blockSize + MEM_HEADER_SIZE
		} else {
			scanLoc++
		}
	}

	newID := highestID + 1
	loc := uint64(0)

	for loc+MEM_HEADER_SIZE+size <= uint64(len(emu.memory)) {
		header := binary.LittleEndian.Uint64(emu.memory[loc : loc+8])
		if header == 0 {
			free_space := uint64(0)
			at := loc
			for at < uint64(len(emu.memory)) {
				if emu.memory[at] != 0 {
					break
				}
				at++
				free_space++
			}

			if free_space >= MEM_HEADER_SIZE+size {
				binary.LittleEndian.PutUint64(emu.memory[loc:loc+8], newID)
				binary.LittleEndian.PutUint64(emu.memory[loc+8:loc+16], size)
				mID = ungo.Some(newID)
				return
			}
			loc += free_space
			continue
		}

		blockSize := binary.LittleEndian.Uint64(emu.memory[loc+8 : loc+16])
		loc += blockSize + MEM_HEADER_SIZE
	}

	emu.memory = append(emu.memory, make([]byte, MEM_HEADER_SIZE+size)...)
	mID = emu.Allocate(size)

	return
}

func (emu *Emulator) Free(mID uint64) (found bool) {
	if emu.mu.TryLock() {
		defer emu.mu.Unlock()
	}

	found = false

	loc := uint64(0)
	for loc+MEM_HEADER_SIZE+1 <= uint64(len(emu.memory)) {
		header_bytes := emu.memory[loc : loc+8]
		header := binary.LittleEndian.Uint64(header_bytes)
		if header == 0 {
			loc++
			continue
		}
		size_bytes := emu.memory[loc+8 : loc+16]
		size := binary.LittleEndian.Uint64(size_bytes)

		if header == mID {
			copy(emu.memory[loc:], make([]byte, size+MEM_HEADER_SIZE))
			found = true
			return
		}

		loc += MEM_HEADER_SIZE + size
	}

	return
}

func (emu *Emulator) Pointer(mID uint64) (location uint64) {
	if emu.mu.TryRLock() {
		defer emu.mu.RUnlock()
	}
	location = 0
	for location+MEM_HEADER_SIZE+1 <= uint64(len(emu.memory)) {
		header_bytes := emu.memory[location : location+8]
		header := binary.LittleEndian.Uint64(header_bytes)
		if header == 0 {
			location++
			continue
		}
		size_bytes := emu.memory[location+8 : location+16]
		size := binary.LittleEndian.Uint64(size_bytes)

		if header == mID {
			location += MEM_HEADER_SIZE
			return
		}

		location += MEM_HEADER_SIZE + size
	}
	return
}

func (emu *Emulator) WriteTo(location uint64, data []byte) {

	if emu.mu.TryLock() {
		defer emu.mu.Unlock()
	}

	copy(emu.memory[location:location+uint64(len(data))], data)
}

func (emu *Emulator) WriteAt(mID uint64, data []byte) (found bool) {
	loc := emu.Pointer(mID)

	if loc == 0 {
		return false
	}

	emu.WriteTo(loc, data)

	return true
}
