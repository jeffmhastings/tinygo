
/* Unused, but here to silence a linker warning. */
ENTRY(Reset_Handler)

/* define output sections */
SECTIONS
{
    .flash_begin : {
        __flash_binary_start = .;
    } > FLASH_TEXT

    .boot2 : {
        __boot2_start__ = .;
        KEEP (*(.boot2))
        __boot2_end__ = .;
    } > FLASH_TEXT

    ASSERT(__boot2_end__ - __boot2_start__ == 256,
        "ERROR: Pico second stage bootloader must be 256 bytes in size")

    /* Program code and read-only data goes to FLASH_TEXT. */
    .text :
    {
        KEEP(*(.isr_vector))
        *(.text)
        *(.text.*)
        *(.rodata)
        *(.rodata.*)
        . = ALIGN(4);
    } >FLASH_TEXT

    .tinygo_stacksizes :
    {
        *(.tinygo_stacksizes)
    } > FLASH_TEXT

    /* Put the stack at the bottom of RAM, so that the application will
     * crash on stack overflow instead of silently corrupting memory.
     * See: http://blog.japaric.io/stack-overflow-protection/ */
    .stack (NOLOAD) :
    {
        . = ALIGN(4);
        . += _stack_size;
        _stack_top = .;
    } >RAM

    /* Start address (in flash) of .data, used by startup code. */
    _sidata = LOADADDR(.data);

    /* Globals with initial value */
    .data :
    {
        . = ALIGN(4);
        _sdata = .;        /* used by startup code */
        *(.data)
        *(.data.*)
        . = ALIGN(4);
        _edata = .;        /* used by startup code */
    } >RAM AT>FLASH_TEXT

    /* Zero-initialized globals  */
    .bss :
    {
        . = ALIGN(4);
        _sbss = .;         /* used by startup code */
        *(.bss)
        *(.bss.*)
        *(COMMON)
        . = ALIGN(4);
        _ebss = .;         /* used by startup code */
    } >RAM

    /DISCARD/ :
    {
        *(.ARM.exidx)      /* causes 'no memory region specified' error in lld */
        *(.ARM.exidx.*)    /* causes spurious 'undefined reference' errors */
    }
}

/* For the memory allocator. */
_heap_start = _ebss;
_heap_end = ORIGIN(RAM) + LENGTH(RAM);
_globals_start = _sdata;
_globals_end = _ebss;
