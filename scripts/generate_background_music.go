package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

const (
	sampleRate = 32000
	bpm        = 82.0
	barCount   = 8
)

type stereoBuffer struct {
	left  []float64
	right []float64
}

type harmony struct {
	root  int
	chord []int
}

func midiFrequency(note int) float64 {
	return 440 * math.Pow(2, float64(note-69)/12)
}

func clamp(value, minimum, maximum float64) float64 {
	return math.Max(minimum, math.Min(maximum, value))
}

func (b stereoBuffer) addTone(start, duration float64, note int, amplitude, pan float64, instrument string) {
	startSample := max(0, int(start*sampleRate))
	endSample := int((start + duration) * sampleRate)
	frequency := midiFrequency(note)
	leftPan := math.Cos(clamp(pan, 0, 1) * math.Pi / 2)
	rightPan := math.Sin(clamp(pan, 0, 1) * math.Pi / 2)

	for sample := startSample; sample < endSample; sample++ {
		t := float64(sample)/sampleRate - start
		targetSample := sample % len(b.left)
		attack := clamp(t/0.045, 0, 1)
		release := clamp((duration-t)/0.42, 0, 1)
		envelope := attack * release
		phase := 2 * math.Pi * frequency * t
		wave := math.Sin(phase)

		switch instrument {
		case "pad":
			envelope = clamp(t/0.65, 0, 1) * clamp((duration-t)/0.9, 0, 1)
			wave = 0.72*math.Sin(phase) + 0.20*math.Sin(phase*2+0.15) + 0.08*math.Sin(phase*0.5)
		case "keys":
			decay := 0.48 + 0.52*math.Exp(-2.4*t)
			envelope *= decay
			wave = 0.76*math.Sin(phase) + 0.18*math.Sin(phase*2) + 0.06*math.Sin(phase*3)
		case "bass":
			envelope = clamp(t/0.03, 0, 1) * clamp((duration-t)/0.28, 0, 1)
			wave = 0.84*math.Sin(phase) + 0.16*math.Sin(phase*2)
		case "bell":
			decay := math.Exp(-2.8 * t)
			envelope *= 0.45 + 0.55*decay
			wave = 0.68*math.Sin(phase) + 0.22*math.Sin(phase*2.01) + 0.10*math.Sin(phase*3.98)
		}

		value := wave * envelope * amplitude
		b.left[targetSample] += value * leftPan
		b.right[targetSample] += value * rightPan
	}
}

func (b stereoBuffer) addSoftPulse(start, amplitude float64) {
	duration := 0.32
	startSample := max(0, int(start*sampleRate))
	endSample := min(len(b.left), int((start+duration)*sampleRate))
	for sample := startSample; sample < endSample; sample++ {
		t := float64(sample)/sampleRate - start
		frequency := 72 - 24*(t/duration)
		wave := math.Sin(2 * math.Pi * frequency * t)
		envelope := math.Exp(-11 * t)
		value := wave * envelope * amplitude
		b.left[sample] += value * 0.707
		b.right[sample] += value * 0.707
	}
}

func (b stereoBuffer) addShaker(start, amplitude, pan float64, seed uint32) {
	duration := 0.11
	startSample := max(0, int(start*sampleRate))
	endSample := min(len(b.left), int((start+duration)*sampleRate))
	leftPan := math.Cos(clamp(pan, 0, 1) * math.Pi / 2)
	rightPan := math.Sin(clamp(pan, 0, 1) * math.Pi / 2)
	previous := 0.0
	for sample := startSample; sample < endSample; sample++ {
		seed = seed*1664525 + 1013904223
		noise := float64(seed>>8)/float64(1<<24)*2 - 1
		highPass := noise - previous*0.82
		previous = noise
		t := float64(sample)/sampleRate - start
		envelope := math.Exp(-34 * t)
		value := highPass * envelope * amplitude
		b.left[sample] += value * leftPan
		b.right[sample] += value * rightPan
	}
}

func buildMusic() stereoBuffer {
	beat := 60 / bpm
	duration := float64(barCount) * 4 * beat
	sampleCount := int(duration * sampleRate)
	buffer := stereoBuffer{left: make([]float64, sampleCount), right: make([]float64, sampleCount)}

	progression := []harmony{
		{48, []int{48, 52, 55, 59}},
		{45, []int{45, 48, 52, 55}},
		{41, []int{41, 45, 48, 52}},
		{43, []int{43, 47, 50, 57}},
		{48, []int{48, 52, 55, 59}},
		{40, []int{40, 43, 47, 50}},
		{41, []int{41, 45, 48, 52}},
		{43, []int{43, 47, 50, 55}},
	}
	melodies := [][]int{
		{64, 67, 69, 67}, {64, 60, 64, 67}, {65, 69, 67, 64}, {62, 67, 71, 69},
		{67, 72, 71, 67}, {64, 67, 71, 67}, {69, 67, 65, 64}, {62, 64, 67, 71},
	}
	arpPattern := []int{0, 2, 1, 3, 2, 1, 0, 2}

	for bar := 0; bar < barCount; bar++ {
		barStart := float64(bar) * 4 * beat
		harmony := progression[bar]
		for chordIndex, note := range harmony.chord {
			pan := 0.28 + float64(chordIndex)*0.15
			// 延长到下一小节并循环叠加，让最后一小节自然过渡回第一小节。
			buffer.addTone(barStart, 5.5*beat, note+12, 0.037, pan, "pad")
		}

		for beatIndex := 0; beatIndex < 4; beatIndex++ {
			time := barStart + float64(beatIndex)*beat
			bassNote := harmony.root
			if beatIndex == 3 {
				bassNote = harmony.chord[2]
			}
			buffer.addTone(time, beat*0.82, bassNote, 0.095, 0.46, "bass")
			if beatIndex%2 == 0 {
				buffer.addSoftPulse(time, 0.075)
			}
			buffer.addTone(time, beat*0.78, melodies[bar][beatIndex], 0.075, 0.38+0.12*float64(beatIndex%2), "bell")
		}

		for eighth := 0; eighth < 8; eighth++ {
			time := barStart + float64(eighth)*beat/2
			note := harmony.chord[arpPattern[eighth]] + 24
			buffer.addTone(time, beat*0.42, note, 0.042, 0.68, "keys")
			if eighth%2 == 1 {
				buffer.addShaker(time, 0.018, 0.72, uint32(1000+bar*17+eighth))
			}
		}
	}

	// 用极短的边界桥接保证首尾采样值连续，AudioBuffer 循环时不会出现咔哒声。
	bridgeSamples := int(0.004 * sampleRate)
	for offset := 0; offset < bridgeSamples; offset++ {
		index := len(buffer.left) - bridgeSamples + offset
		blend := float64(offset+1) / float64(bridgeSamples)
		buffer.left[index] = buffer.left[index]*(1-blend) + buffer.left[0]*blend
		buffer.right[index] = buffer.right[index]*(1-blend) + buffer.right[0]*blend
	}

	peak := 0.0
	for sample := range buffer.left {
		peak = math.Max(peak, math.Abs(buffer.left[sample]))
		peak = math.Max(peak, math.Abs(buffer.right[sample]))
	}
	if peak > 0 {
		scale := 0.82 / peak
		for sample := range buffer.left {
			buffer.left[sample] *= scale
			buffer.right[sample] *= scale
		}
	}
	return buffer
}

func writeWAV(path string, audio stereoBuffer) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dataSize := uint32(len(audio.left) * 2 * 2)
	byteRate := uint32(sampleRate * 2 * 2)
	blockAlign := uint16(2 * 2)
	write := func(value any) error { return binary.Write(file, binary.LittleEndian, value) }
	if _, err := file.WriteString("RIFF"); err != nil {
		return err
	}
	if err := write(uint32(36) + dataSize); err != nil {
		return err
	}
	if _, err := file.WriteString("WAVEfmt "); err != nil {
		return err
	}
	for _, value := range []any{uint32(16), uint16(1), uint16(2), uint32(sampleRate), byteRate, blockAlign, uint16(16)} {
		if err := write(value); err != nil {
			return err
		}
	}
	if _, err := file.WriteString("data"); err != nil {
		return err
	}
	if err := write(dataSize); err != nil {
		return err
	}
	for sample := range audio.left {
		left := int16(clamp(audio.left[sample], -1, 1) * math.MaxInt16)
		right := int16(clamp(audio.right[sample], -1, 1) * math.MaxInt16)
		if err := write(left); err != nil {
			return err
		}
		if err := write(right); err != nil {
			return err
		}
	}
	return file.Sync()
}

func main() {
	output := flag.String("output", "frontend/public/audio/lifegame-theme.wav", "output WAV path")
	flag.Parse()
	audio := buildMusic()
	if err := writeWAV(*output, audio); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("generated %s (%d stereo samples at %d Hz)\n", *output, len(audio.left), sampleRate)
}
