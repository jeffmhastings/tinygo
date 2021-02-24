// +build rp2040

package machine

import (
	"device/rp2040"
)

type PinMode uint8

const (
	PinOutput PinMode = 0
)

// Set the pin to high or low.
// Warning: only use this on an output pin!
func (p Pin) Set(high bool) {
	//pin is unit8 or uint32?
	if high {
		rp2040.SIO.GPIO_OUT.SetBits(1 << p)
	} else {
		rp2040.SIO.GPIO_OUT.ClearBits(1 << p)
	}
}

// Configure this pin with the given I/O settings.
func (p Pin) Configure(config PinConfig) {
	// Clear Output Enable bit
	rp2040.SIO.GPIO_OE_CLR.SetBits(1 << 25)

	// Clear output output value
	rp2040.SIO.GPIO_OUT_CLR.SetBits(1 << 25)

	// Set Input enable, Clear output disable
	rp2040.PADS_BANK0.GPIO25.ReplaceBits(rp2040.PADS_BANK0_GPIO0_IE_Msk, rp2040.PADS_BANK0_GPIO0_IE_Pos|rp2040.PADS_BANK0_GPIO0_OD_Pos, 0)

	// Select function SIO
	rp2040.IO_BANK0.GPIO25_CTRL.Set(rp2040.IO_BANK0_GPIO0_CTRL_FUNCSEL_SIO_0 << rp2040.IO_BANK0_GPIO0_CTRL_FUNCSEL_Pos)

	// Select output enable
	rp2040.SIO.GPIO_OE_SET.SetBits(1 << 25)
}
