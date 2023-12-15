package robot

import (
	. "machine"
	"time"

	"github.com/sparques/pwm"
)

type Servo struct {
	pin Pin
	pwm pwm.Group
}

func NewServo(pin Pin) *Servo {
	pin.Configure(PinConfig{Mode: PinPWM})
	p := pwm.Get(pin)
	p.Configure(PWMConfig{Period: uint64(20 * time.Millisecond)})

	return &Servo{
		pin: pin,
		pwm: p,
	}
}

func (s *Servo) SetFloat(vector float32) {
	sch, _ := s.pwm.Channel(s.pin)
	ns := float32(s.pwm.Top() / 20000)
	s.pwm.Set(sch, uint32(ns*(1500+1000*vector)))
}

// Set sets the servo, the nominal range is -500 to 500, but some servos respond
// outside of this range.
func (s *Servo) Set(vector int32) {
	sch, _ := s.pwm.Channel(s.pin)
	ns := (int64(s.pwm.Top()) * (int64(vector) + 1500)) / 20000
	s.pwm.Set(sch, uint32(ns))
}
