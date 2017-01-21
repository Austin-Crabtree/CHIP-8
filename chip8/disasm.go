/* Copyright (c) 2017 Jeffrey Massung
 *
 * This software is provided 'as-is', without any express or implied
 * warranty.  In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would be
 *    appreciated but is not required.
 *
 * 2. Altered source versions must be plainly marked as such, and must not be
 *    misrepresented as being the original software.
 *
 * 3. This notice may not be removed or altered from any source distribution.
 */

package chip8

import "fmt"

/// Disassemble a CHIP-8 instruction.
///
func (vm *CHIP_8) Disassemble(i uint) string {
	if int(i) >= len(vm.Memory) - 1 {
		return ""
	}

	// fetch the instruction at this location
	inst := uint(vm.Memory[i])<<8 | uint(vm.Memory[i+1])

	// end of program memory?
	if inst == 0 {
		return fmt.Sprintf("%04X -", i)
	}

	// 12-bit literal Address
	a := inst & 0xFFF

	// byte and nibble literals
	b := byte(inst & 0xFF)
	n := byte(inst & 0xF)

	// vx and vy registers
	x := inst >> 8 & 0xF
	y := inst >> 4 & 0xF

	// instruction decoding
	if inst == 0x00E0 {
		return fmt.Sprintf("%04X - CLS", i)
	} else if inst == 0x00EE {
		return fmt.Sprintf("%04X - RET", i)
	} else if inst == 0x00FE {
		return fmt.Sprintf("%04X - LOW", i)
	} else if inst == 0x00FF {
		return fmt.Sprintf("%04X - HIGH", i)
	} else if inst == 0x00FB {
		return fmt.Sprintf("%04X - SCR", i)
	} else if inst == 0x00FC {
		return fmt.Sprintf("%04X - SCL", i)
	} else if inst == 0x00FD {
		return fmt.Sprintf("%04X - EXIT", i)
	} else if inst&0xFFF0 == 0x00B0 {
		return fmt.Sprintf("%04X - SCU    %d", i, n)
	} else if inst&0xFFF0 == 0x00C0 {
		return fmt.Sprintf("%04X - SCD    %d", i, n)
	} else if inst&0xF000 == 0x0000 {
		return fmt.Sprintf("%04X - SYS    #%04X", i, a)
	} else if inst&0xF000 == 0x1000 {
		return fmt.Sprintf("%04X - JP     #%04X", i, a)
	} else if inst&0xF000 == 0x2000 {
		return fmt.Sprintf("%04X - CALL   #%04X", i, a)
	} else if inst&0xF000 == 0x3000 {
		return fmt.Sprintf("%04X - SE     V%X, #%02X", i, x, b)
	} else if inst&0xF000 == 0x4000 {
		return fmt.Sprintf("%04X - SNE    V%X, #%02X", i, x, b)
	} else if inst&0xF00F == 0x5000 {
		return fmt.Sprintf("%04X - SE     V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x5001 {
		return fmt.Sprintf("%04X - SGT    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x5002 {
		return fmt.Sprintf("%04X - SLT    V%X, V%X", i, x, y)
	} else if inst&0xF000 == 0x6000 {
		return fmt.Sprintf("%04X - LD     V%X, #%02X", i, x, b)
	} else if inst&0xF000 == 0x7000 {
		return fmt.Sprintf("%04X - ADD    V%X, #%02X", i, x, b)
	} else if inst&0xF00F == 0x8000 {
		return fmt.Sprintf("%04X - LD     V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x8001 {
		return fmt.Sprintf("%04X - OR     V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x8002 {
		return fmt.Sprintf("%04X - AND    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x8003 {
		return fmt.Sprintf("%04X - XOR    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x8004 {
		return fmt.Sprintf("%04X - ADD    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x8005 {
		return fmt.Sprintf("%04X - SUB    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x8006 {
		return fmt.Sprintf("%04X - SHR    V%X", i, x)
	} else if inst&0xF00F == 0x8007 {
		return fmt.Sprintf("%04X - SUBN   V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x800E {
		return fmt.Sprintf("%04X - SHL    V%X", i, x)
	} else if inst&0xF00F == 0x9000 {
		return fmt.Sprintf("%04X - SNE    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x9001 {
		return fmt.Sprintf("%04X - MUL    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x9002 {
		return fmt.Sprintf("%04X - DIV    V%X, V%X", i, x, y)
	} else if inst&0xF00F == 0x9003 {
		return fmt.Sprintf("%04X - BCD    V%X, V%X", i, x, y)
	} else if inst&0xF000 == 0xA000 {
		return fmt.Sprintf("%04X - LD     I, #%04X", i, a)
	} else if inst&0xF000 == 0xB000 {
		return fmt.Sprintf("%04X - JP     V0, #%04X", i, a)
	} else if inst&0xF000 == 0xC000 {
		return fmt.Sprintf("%04X - RND    V%X, #%02X", i, x, b)
	} else if inst&0xF000 == 0xD000 {
		return fmt.Sprintf("%04X - DRW    V%X, V%X, %d", i, x, y, n)
	} else if inst&0xF0FF == 0xE09E {
		return fmt.Sprintf("%04X - SKP    V%X", i, x)
	} else if inst&0xF0FF == 0xE0A1 {
		return fmt.Sprintf("%04X - SKNP   V%X", i, x)
	} else if inst&0xF0FF == 0xF007 {
		return fmt.Sprintf("%04X - LD     V%X, DT", i, x)
	} else if inst&0xF0FF == 0xF00A {
		return fmt.Sprintf("%04X - LD     V%X, K", i, x)
	} else if inst&0xF0FF == 0xF015 {
		return fmt.Sprintf("%04X - LD     DT, V%X", i, x)
	} else if inst&0xF0FF == 0xF018 {
		return fmt.Sprintf("%04X - LD     ST, V%X", i, x)
	} else if inst&0xF0FF == 0xF01E {
		return fmt.Sprintf("%04X - ADD    I, V%X", i, x)
	} else if inst&0xF0FF == 0xF029 {
		return fmt.Sprintf("%04X - LD     F, V%X", i, x)
	} else if inst&0xF0FF == 0xF030 {
		return fmt.Sprintf("%04X - LD     HF, V%X", i, x)
	} else if inst&0xF0FF == 0xF033 {
		return fmt.Sprintf("%04X - BCD    V%X", i, x)
	} else if inst&0xF0FF == 0xF055 {
		return fmt.Sprintf("%04X - LD     [I], V%X", i, x)
	} else if inst&0xF0FF == 0xF065 {
		return fmt.Sprintf("%04X - LD     V%X, [I]", i, x)
	} else if inst&0xF0FF == 0xF075 {
		return fmt.Sprintf("%04X - LD     R, V%X", i, x)
	} else if inst&0xF0FF == 0xF085 {
		return fmt.Sprintf("%04X - LD     V%X, R", i, x)
	} else if inst&0xF0FF == 0xF094 {
		return fmt.Sprintf("%04X - LD     A, V%X", i, x)
	}

	// unknown instruction
	return fmt.Sprintf("%04X - ??", i)
}