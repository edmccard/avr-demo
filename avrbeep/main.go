package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"fmt"
	"github.com/edmccard/avr-sim/core"
	"github.com/edmccard/avr-sim/instr"
	"strings"
	"time"
)

var cycleCount uint64

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	cpu := core.Cpu{}
	mem := NewDemoMem(&cpu, strings.NewReader(program))
	decoder := instr.NewDecoder(instr.NewSetEnhanced8k())

	ticker := time.NewTicker(20 * time.Millisecond)
	quit := make(chan struct{})

	spk, err := NewSpeaker(1000000)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	mem.outports[0x38] = spk.write

	go func() {
		cycles := uint(0)
		elapsed := uint(0)
		started := false

		for {
			select {
			case <-ticker.C:
				for cycles < 20000 {
					elapsed = cpu.Step(mem, &decoder)
					cycles += elapsed
					cycleCount += uint64(elapsed)
				}
				cycles -= 20000
				if !started {
					err = spk.Start()
					if err != nil {
						fmt.Println("ERROR:", err)
						return
					}
					started = true
				}
			case <-quit:
				ticker.Stop()
				spk.stream.Stop()
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)
	close(quit)
}

// BEEP.ASM from https://sites.google.com/site/avrasmintro/
const program = `
:020000020000FC
:100000000FE50DBF04E00EBF0FEF07BB55270FEF45
:1000100008BB06D0002708BB03D05A95C1F7FFCF15
:0A002000662700006A95E9F70895CD
:00000001FF
`
