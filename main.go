package main

func main() {
	emu := NewEmulator()

	programMID := emu.InitProgram([]byte{
		0, 0, 0, 0, 0, 0, 0, 0, // Valid NO OP
	})

	if programMID.HasValue() {
		defer emu.Free(programMID.Value())
	}
}
