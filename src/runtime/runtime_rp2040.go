// +build rp2040

package runtime

import (
	"device/arm"
	"device/rp2040"
)

type timeUnit int64

const (
	asyncScheduler = false

	MHZ = 10 ^ 6

	// RP2040 reference design uses a 12 MHz crystal
	// rp2040-datasheet.pdf 2.16.1 pg. 234
	XOSC_MHZ = uint32(12)

	// rp2040-datasheet.pdf 2.16.3 Startup Delay pg. 235
	XOSC_STARTUP_DELAY = (((XOSC_MHZ * 1000000) / 1000) + 128) / 256
)

//export Reset_Handler
func main() {
	preinit()
	run()
	abort()
}

func init() {
	initReset()
	initClocks()
}

func initReset() {
	// Reset all peripherals to put system into a known state,
	// - except for QSPI pads and the XIP IO bank, as this is fatal if running from flash
	// - and the PLLs, as this is fatal if clock muxing has not been reset on this boot
	rp2040.RESETS.RESET.SetBits(rp2040.RESETS_RESET_USBCTRL |
		rp2040.RESETS_RESET_UART1 |
		rp2040.RESETS_RESET_UART0 |
		rp2040.RESETS_RESET_TIMER |
		rp2040.RESETS_RESET_TBMAN |
		rp2040.RESETS_RESET_SYSINFO |
		rp2040.RESETS_RESET_SYSCFG |
		rp2040.RESETS_RESET_SPI1 |
		rp2040.RESETS_RESET_SPI0 |
		rp2040.RESETS_RESET_RTC |
		rp2040.RESETS_RESET_PWM |
		rp2040.RESETS_RESET_PIO1 |
		rp2040.RESETS_RESET_PIO0 |
		rp2040.RESETS_RESET_PADS_BANK0 |
		rp2040.RESETS_RESET_JTAG |
		rp2040.RESETS_RESET_IO_BANK0 |
		rp2040.RESETS_RESET_I2C1 |
		rp2040.RESETS_RESET_I2C0 |
		rp2040.RESETS_RESET_DMA |
		rp2040.RESETS_RESET_BUSCTRL |
		rp2040.RESETS_RESET_ADC)

	// Remove reset from peripherals which are clocked only by clk_sys and
	// clk_ref. Other peripherals stay in reset until we've configured clocks.
	peripheralsToDeassert := uint32(rp2040.RESETS_RESET_BUSCTRL |
		rp2040.RESETS_RESET_DMA |
		rp2040.RESETS_RESET_I2C0 |
		rp2040.RESETS_RESET_I2C1 |
		rp2040.RESETS_RESET_IO_BANK0 |
		rp2040.RESETS_RESET_JTAG |
		rp2040.RESETS_RESET_PADS_BANK0 |
		rp2040.RESETS_RESET_PIO0 |
		rp2040.RESETS_RESET_PIO1 |
		rp2040.RESETS_RESET_PWM |
		rp2040.RESETS_RESET_SYSCFG |
		rp2040.RESETS_RESET_SYSINFO |
		rp2040.RESETS_RESET_TBMAN |
		rp2040.RESETS_RESET_TIMER)

	rp2040.RESETS.RESET.ClearBits(peripheralsToDeassert)
	for !rp2040.RESETS.RESET_DONE.HasBits(peripheralsToDeassert) {
	}
}

func initClocks() {
	// Start the watchdog tick
	rp2040.WATCHDOG.TICK.Set(XOSC_MHZ | rp2040.WATCHDOG_TICK_ENABLE)

	// Reset clock resus
	rp2040.CLOCKS.CLK_SYS_RESUS_CTRL.Set(0)

	// Assumes 1-15 MHz input
	rp2040.XOSC.CTRL.Set(rp2040.XOSC_CTRL_FREQ_RANGE_1_15MHZ)

	// Set xosc startup delay
	rp2040.XOSC.STARTUP.Set(XOSC_STARTUP_DELAY)

	// Set the enable bit now that we have set freq range and startup delay
	rp2040.XOSC.CTRL.SetBits(rp2040.XOSC_CTRL_ENABLE_ENABLE << rp2040.XOSC_CTRL_ENABLE_Pos)

	// Wait for XOSC to be stable
	for !rp2040.XOSC.STATUS.HasBits(rp2040.XOSC_STATUS_STABLE) {
	}
}

func postinit() {}

func ticksToNanoseconds(ticks timeUnit) int64 {
	return int64(ticks) * 1000
}

func nanosecondsToTicks(ns int64) timeUnit {
	return timeUnit(ns / 1000)
}

// number of ticks (microseconds) since start.
func ticks() timeUnit {
	// rp2040-datasheet.pdf 4.6.2 Counter pg. 557

	// Read from the lower register first to latch the higher
	lr := rp2040.TIMER.TIMELR.Get()

	// Read from the higher register
	hr := rp2040.TIMER.TIMEHR.Get()

	return timeUnit(hr<<32 | lr)
}

// sleepTicks should sleep for specific number of microseconds.
func sleepTicks(d timeUnit) {
	start := rp2040.TIMER.TIMERAWL.Get()
	for rp2040.TIMER.TIMERAWL.Get()-start < uint32(d) {
	}
}

func putchar(c byte) {
}

func waitForEvents() {
	arm.Asm("wfe")
}
