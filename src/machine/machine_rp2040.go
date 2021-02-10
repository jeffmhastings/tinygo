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
	// Set input enable on, output disable off
	//    hw_write_masked(&padsbank0_hw->io[gpio],
	//                   PADS_BANK0_GPIO0_IE_BITS,
	//                   PADS_BANK0_GPIO0_IE_BITS | PADS_BANK0_GPIO0_OD_BITS
	//    );
	//    // Zero all fields apart from fsel; we want this IO to do what the peripheral tells it.
	//    // This doesn't affect e.g. pullup/pulldown, as these are in pad controls.
	//    iobank0_hw->io[gpio].ctrl = fn << IO_BANK0_GPIO0_CTRL_FUNCSEL_LSB;

	switch p {
	case 25:
		rp2040.IO_BANK0.GPIO25_CTRL.Set(rp2040.IO_BANK0_GPIO0_CTRL_FUNCSEL_SIO_0 << rp2040.IO_BANK0_GPIO0_CTRL_FUNCSEL_Pos)
	}
}

//---------- UART related types and code

// UART representation
//type UART struct {
//	Buffer    *RingBuffer
//	Bus       *rp2040.UART0_Type
//	Interrupt interrupt.Interrupt
//}

// Configure the TX and RX pins
//func (uart UART) configurePins(config UARTConfig) {

// pins
//switch config.TX {
//case UART_ALT_TX_PIN:
//	// use alternate TX/RX pins via AFIO mapping
//	stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_AFIOEN)
//	if uart.Bus == stm32.USART1 {
//		stm32.AFIO.MAPR.SetBits(stm32.AFIO_MAPR_USART1_REMAP)
//	} else if uart.Bus == stm32.USART2 {
//		stm32.AFIO.MAPR.SetBits(stm32.AFIO_MAPR_USART2_REMAP)
//	}
//default:
//	// use standard TX/RX pins PA9 and PA10
//}
//config.TX.Configure(PinConfig{Mode: PinOutput50MHz + PinOutputModeAltPushPull})
//config.RX.Configure(PinConfig{Mode: PinInputModeFloating})
//}
