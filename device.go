package main

import "github.com/loeredami/ungo"

func (emu *Emulator) NewPID() (PID uint64) {
	PID = 0

	for _, w := range emu.workers {
		if PID <= w.PID {
			PID = w.PID + 1
		}
	}

	return
}

func (emu *Emulator) WorkerSpawnAt(mID uint64) (success bool) {
	success = false
	loc := emu.Pointer(mID)

	if loc == 0 {
		return
	}

	emu.mu.Lock()
	defer emu.mu.Unlock()
	w := &Worker{
		PID:      emu.NewPID(),
		IP:       loc,
		emulator: emu,
		stack:    ungo.NewStack[uint64](),
		runnin:   true,
	}

	defer func() { go w.Work() }()

	emu.workers = append(emu.workers, w)

	success = true
	return
}
