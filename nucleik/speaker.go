package main

import (
	"code.google.com/p/portaudio-go/portaudio"
	"github.com/edmccard/avr-sim/core"
)

type Speaker struct {
	stream       *portaudio.Stream
	channel      chan float32
	curSample    float32
	avgBuf       []float32
	avgIdx       int
	cycPerSample uint
	pin          byte
	lastToggle   uint64
	dbgCount     int
}

func NewSpeaker(hertz uint) (*Speaker, error) {
	spk := &Speaker{curSample: -1.0}
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	parameters.Output.Channels = 1
	parameters.SampleRate = 44100
	stream, err := portaudio.OpenStream(parameters, spk.Callback)
	if err != nil {
		return nil, err
	}
	spk.cycPerSample = hertz / uint(parameters.SampleRate)
	spk.channel = make(chan float32, 8192)
	spk.avgBuf = make([]float32, spk.cycPerSample)
	spk.stream = stream
	return spk, nil
}

func (spk *Speaker) Start() error {
	return spk.stream.Start()
}

func (spk *Speaker) Stop() error {
	return spk.stream.Stop()
}

func (spk *Speaker) Callback(out []float32) {
	for i := range out {
		select {
		case sample := <-spk.channel:
			out[i] = sample
		default:
			out[i] = 0
		}
	}
	spk.dbgCount = 0
}

func (spk *Speaker) write(addr core.Addr, val byte) {
	val &= 1
	if val == spk.pin {
		return
	}
	spk.pin = val
	// TODO: cycleCount is a global variable. Fix?
	spk.makeSamples(cycleCount)
	spk.curSample *= -1.0
}

func (spk *Speaker) makeSamples(curCycle uint64) {
	elapsed := curCycle - spk.lastToggle
	spk.lastToggle = curCycle

	if spk.avgIdx != 0 {
		for ; spk.avgIdx < len(spk.avgBuf); spk.avgIdx++ {
			if elapsed == 0 {
				break
			}
			spk.avgBuf[spk.avgIdx] = spk.curSample
			elapsed--
		}
		if spk.avgIdx == len(spk.avgBuf) {
			spk.avgIdx = 0
			avg := float32(0.0)
			for _, sample := range spk.avgBuf {
				avg += sample
			}
			avg /= float32(len(spk.avgBuf))
			spk.sendSample(avg)
		}
	}

	for i := uint64(0); i < elapsed/uint64(spk.cycPerSample); i++ {
		spk.sendSample(spk.curSample)
	}

	avgCycs := elapsed % uint64(spk.cycPerSample)
	if avgCycs != 0 {
		for spk.avgIdx = 0; spk.avgIdx < int(avgCycs); spk.avgIdx++ {
			spk.avgBuf[spk.avgIdx] = spk.curSample
		}
	}
}

func (spk *Speaker) sendSample(sample float32) {
	select {
	case spk.channel <- sample:
		spk.dbgCount++
	default:
	}
}
