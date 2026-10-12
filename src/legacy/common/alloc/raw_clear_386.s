#include "textflag.h"

// Clear exactly n bytes of caller-owned unmanaged memory.
TEXT ·rawClear(SB), NOSPLIT, $0-8
    MOVL ptr+0(FP), DI
    MOVL n+4(FP), CX
    XORL AX, AX
    REP
    STOSB
    RET
