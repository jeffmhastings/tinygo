// +build rp2040

package runtime

import (
	"device/arm"
	"device/rp2040"
	//"runtime/interrupt"
)

type timeUnit int64

const asyncScheduler = false
const XOSC_MHZ = uint32(12)

func init() {
	//machine.UART0.Configure(machine.UARTConfig{})

	// Start the watchdog tick
	rp2040.WATCHDOG.TICK.Set(XOSC_MHZ | rp2040.WATCHDOG_TICK_ENABLE)
}

func ticksToNanoseconds(ticks timeUnit) int64 {
	return int64(ticks) * 1000
}

func nanosecondsToTicks(ns int64) timeUnit {
	return timeUnit(ns / 1000)
}

// sleepTicks should sleep for specific number of microseconds.
func sleepTicks(d timeUnit) {
	rp2040.TIMER.INTE.SetBits(1 << 0)
	rp2040.SYSCFG.PROC0_NMI_MASK.SetBits(1 << rp2040.IRQ_TIMER_IRQ_0)
	rp2040.TIMER.ALARM0.Set(rp2040.TIMER.TIMERAWL.Get() + uint32(d))

	// Check for armed bit to be set to 0
	for rp2040.TIMER.ARMED.HasBits(1 << 0) {
		arm.Asm("wfi")
	}
	// clear the latched interrupt
	rp2040.TIMER.INTE.SetBits(1 << 0)
}

// number of ticks (microseconds) since start.
func ticks() timeUnit {

	// Read microseconds from the counter

	// Read from the lower register first to latch the higher
	lr := rp2040.TIMER.TIMELR.Get()

	// Read from the higher register
	hr := rp2040.TIMER.TIMEHR.Get()

	return timeUnit(hr<<32 | lr)
}

func postinit() {}

func putchar(c byte) {
	//machine.UART0.WriteByte(c)
}

//export Reset_Handler
func main() {
	run()
}

func waitForEvents() {
	arm.Asm("wfe")
}
