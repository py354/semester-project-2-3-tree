package utils

import "time"

// Timer - стуктура для измерения времени с возможностью остановки и запуска
type Timer struct {
	start    time.Time
	passed   time.Duration
	isActive bool
}

func NewTimer() *Timer {
	return &Timer{
		start:    time.Now(),
		passed:   0,
		isActive: false,
	}
}

func (t *Timer) Passed() time.Duration {
	if t.isActive {
		return time.Since(t.start) + t.passed
	} else {
		return t.passed
	}
}

func (t *Timer) Start() {
	if t.isActive {
		panic("Timer already active")
	}

	t.isActive = true
	t.start = time.Now()
}

func (t *Timer) Stop() {
	if !t.isActive {
		panic("Timer already stopped")
	}

	t.isActive = false
	t.passed += time.Since(t.start)
}
