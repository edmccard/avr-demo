package main

import (
	"fmt"
	"strings"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/edmccard/avr-sim/atmega8"
	"github.com/edmccard/avr-sim/dev"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	sys := atmega8.NewSystem()
	sys.LoadProgHex(strings.NewReader(program))

	spk, err := dev.NewSpeaker(sys.Timer, 1000000, 0)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	sys.Memory.SetWriter(0x38, spk.Write)

	quit := sys.Go(1000000, 50, spk.OnSlice)
	time.Sleep(1 * time.Second)
	close(quit)
	spk.Stop()
}

// BEEP.ASM from https://sites.google.com/site/avrasmintro/
const program = `
:020000020000FC
:100000000FE50DBF04E00EBF0FEF07BB55270FEF45
:1000100008BB06D0002708BB03D05A95C1F7FFCF15
:0A002000662700006A95E9F70895CD
:00000001FF
`
