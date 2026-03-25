package main

func main() {
	emu := NewEmulator()

	programMID := emu.InitProgram([]byte{
		0, 0, 0, 0, 0, 0, 0, 0, // Valid NO OP
	})

	if programMID.HasValue() {
		defer emu.Free(programMID.Value())
	}

	programMID.IfPresent(func(mID uint64) {
		if !emu.WorkerSpawnAt(mID) {
			panic("Could not spawn main thread")
		}
	})
}
