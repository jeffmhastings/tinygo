
MEMORY
{
    FLASH_TEXT (rx) : ORIGIN = 0x10000000, LENGTH = 2048k
    RAM (rwx)       : ORIGIN = 0x20000000, LENGTH = 256k
}

_stack_size = 2K;

INCLUDE "targets/rp2040arm.ld"
