package robot

import (
	. "machine"

	"github.com/sparques/pwm"
)

// Motor lets you control a brushed DC motor (likely through an H-bridge) via PWM.
// Freewheeling (leaving control pins floating) isn't supported--it will likely
// depend on the H-Bridge you're using.
type Motor struct {
	pins [2]Pin
	pwm  pmw.Group
}

// NewMotor returns a *Motor controlled by the given pins. These pins must
// be on the same PWM slice--it's assumed pins[1] is the B channel of pins[0].
func NewMotor(pins [2]Pin) *Motor {
	pins[0].Configure(PinConfig{Mode: PinPWM})
	pins[1].Configure(PinConfig{Mode: PinPWM})
	p := pwm.Get(pins[0])
	p.Configure(PWMConfig{Period: 1e9 / 20000}) // 20kHz is reasonable

	return &Motor{
		pins: pins,
		pwm:  p,
	}
}

// SetFloat sets the speed of the motor using a float value. -1 is full reverse,
// 0 is stopped and 1 is full forward.
func (m *Motor) SetFloat(vector float32) {
	signal := m.pins[0]
	gnd := m.pins[1]
	if vector < 0 {
		vector = -vector
		signal, gnd = gnd, signal
	}
	if vector > 1 {
		vector = 1
	}

	// always set gnd first, so we don't risk shorting the h-bridge
	gch, _ := m.pwm.Channel(gnd)
	m.pwm.Set(gch, 0)

	sch, _ := m.pwm.Channel(signal)
	m.pwm.Set(sch, uint32(float32(m.pwm.Top())*vector))
}

// Set sets the motor's speed ranging from -32768 through +32767
// Where -32768 is full reverse, 0 is stopped, and +32787 is full forward.
func (m *Motor) Set(vector int16) {
	signal := m.pins[0]
	gnd := m.pins[1]
	if vector < 0 {
		vector = -vector
		signal, gnd = gnd, signal
	}

	gch, _ := m.pwm.Channel(gnd)
	m.pwm.Set(gch, 0)

	sch, _ := m.pwm.Channel(signal)
	m.pwm.Set(sch, uint32(uint64(m.pwm.Top())*uint64(vector)/32767))
}
