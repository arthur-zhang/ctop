package main

import (
	"fmt"
	"time"
)

const (
	FSHIFT  = 11
	FIXED_1 = 1 << FSHIFT // 2048
	EXP_1   = 1884        /* 1/exp(5sec/1min) as fixed-point */
	EXP_5   = 2014        /* 1/exp(5sec/5min) */
	EXP_15  = 2037        /* 1/exp(5sec/15min) */
)

func loadInt(x uint64) uint64 {
	return x >> FSHIFT
}

func loadFrac(x uint64) uint64 {
	return loadInt(((x) & (FIXED_1 - 1)) * 100)
}

type LoadAverage struct {
	avenrun [3]uint64
}

func (l *LoadAverage) Run() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			l.refreshLoad()
		}
	}
}

func calcLoad(load0 uint64, exp uint64, active uint64) uint64 {
	newload := load0*exp + active*(FIXED_1-exp)
	if active >= load0 {
		newload += FIXED_1 - 1
	}

	return newload / FIXED_1
}

func (l *LoadAverage) refreshLoad() {
	tc := GetTaskCount()
	active := tc.Running + tc.Uninterruptible
	if active > 0 {
		active = active * FIXED_1
	}
	l.avenrun[0] = calcLoad(l.avenrun[0], EXP_1, uint64(active))
	l.avenrun[1] = calcLoad(l.avenrun[1], EXP_5, uint64(active))
	l.avenrun[2] = calcLoad(l.avenrun[2], EXP_15, uint64(active))
}

func (l *LoadAverage) GetLoad() string {
	a := l.avenrun[0] + (FIXED_1 / 200) // 向上取整
	b := l.avenrun[1] + (FIXED_1 / 200)
	c := l.avenrun[2] + (FIXED_1 / 200)
	return fmt.Sprintf("%d.%02d, %d.%02d, %d.%02d", loadInt(a), loadFrac(a), loadInt(b), loadFrac(b), loadInt(c), loadFrac(c))
}

func NewLoadAverage() *LoadAverage {
	l := LoadAverage{}
	return &l
}
