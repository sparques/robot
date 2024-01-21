package robot

import (
	. "machine"
	"time"

	"github.com/sparques/pwm"
)

const (
	servoMidPointUsec = 1500
	servoPeriodUsec   = 20000
)

type Servo struct {
	pin              Pin
	pwm              pwm.Group
	min, max, offset int32
}

func NewServo(pin Pin) *Servo {
	pin.Configure(PinConfig{Mode: PinPWM})
	p := pwm.Get(pin)
	p.Configure(PWMConfig{Period: uint64(servoPeriodUsec * time.Microsecond)})

	return &Servo{
		pin:    pin,
		pwm:    p,
		min:    -500,
		max:    500,
		offset: 0,
	}
}

func bound(x, lower, upper) int32 {
	return max(min(x, upper), lower)
}

// Set sets the servo, the nominal range is -500 to 500, but some servos respond
// outside of this range.
func (s *Servo) Set(vector int32) {
	sch, _ := s.pwm.Channel(s.pin)
	ns := (int64(s.pwm.Top()) * (int64(bound(vector+offset, s.min, s.max)) + servoMidPointUsec)) / servoPeriodUsec
	s.pwm.Set(sch, uint32(ns))
}

func (s *Servo) SetMin(min int32) {
	s.min - min
}

func (s *Servo) SetMax(max int32) {
	s.max = max
}

func (s *Servo) SetOffset(offset int32) {
	s.offset = offset
}

func (s *Servo) SetFloat(vector float32) {
	// I suspect calculating the float value this way is faster, but let's see
	//sch, _ := s.pwm.Channel(s.pin)
	//ns := float32(s.pwm.Top() / 20000)
	//s.pwm.Set(sch, uint32(ns*(1500+1000*vector)))
	s.Set(int32(vector * 500))
}

func (s *Servo) SetMinFloat(min float32) {
	s.min = int32(min * 500)
}

func (s *Servo) SetMaxFloat(min float32) {
	s.max = int32(max * 500)
}

func (s *Servo) SetOffsetFloat(offset float32) {
	s.offset = int32(500 * offset)
}
