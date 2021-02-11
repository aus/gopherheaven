//func heaven(shellcodeAddr uintptr, functionAddr uint64, arg0 uint64, arg1 uint64, arg2 uint64, arg3 uint64, extSize uint64, extArgs uint64, pReturn uint32)
TEXT ·heaven(SB),4,$64-64

  MOVL functionAddrA+4(FP), BX
  MOVL BX, 0(SP)
  MOVL functionAddrB+8(FP), BX
  MOVL BX, 4(SP)

  MOVL arg0A+12(FP), BX
  MOVL BX, 8(SP)
  MOVL arg0B+16(FP), BX
  MOVL BX, 12(SP)

  MOVL arg1A+20(FP), BX
  MOVL BX, 16(SP)
  MOVL arg1B+24(FP), BX
  MOVL BX, 20(SP)

  MOVL arg2A+28(FP), BX
  MOVL BX, 24(SP)
  MOVL arg2B+32(FP), BX
  MOVL BX, 28(SP)

  MOVL arg3A+36(FP), BX
  MOVL BX, 32(SP)
  MOVL arg3B+40(FP), BX
  MOVL BX, 36(SP)

  MOVL extSizeA+44(FP), BX
  MOVL BX, 40(SP)
  MOVL extSizeB+48(FP), BX
  MOVL BX, 44(SP)

  MOVL extArgsA+52(FP), BX
  MOVL BX, 48(SP)
  MOVL extArgsB+56(FP), BX
  MOVL BX, 52(SP)

  MOVL ret+60(FP), BX
  MOVL BX, 56(SP)

  MOVL shellcodeAddr+0(FP), AX
  CALL AX

  MOVL AX, ret+72(FP)
  RET
