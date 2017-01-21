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

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

/// Type for scanned tokens.
///
type tokenType uint

/// Lexical assembly tokens.
///
const (
	TOKEN_END tokenType = iota
	TOKEN_CHAR
	TOKEN_LABEL
	TOKEN_ID
	TOKEN_INSTRUCTION
	TOKEN_OPERAND
	TOKEN_V
	TOKEN_R
	TOKEN_I
	TOKEN_EFFECTIVE_ADDRESS
	TOKEN_F
	TOKEN_HF
	TOKEN_K
	TOKEN_DT
	TOKEN_ST
	TOKEN_LIT
	TOKEN_TEXT
	TOKEN_BREAK
	TOKEN_ASSERT
	TOKEN_EQU
	TOKEN_VAR
	TOKEN_HERE
	TOKEN_SUPER
	TOKEN_EXTENDED
	TOKEN_ASCII
)

/// A parsed, lexical token.
///
type token struct {
	typ tokenType

	// tokens can have an optional value associated with them
	val interface{}
}

/// CHIP-8 assembler token scanner.
///
type tokenScanner struct {
	bytes []byte

	// scan position
	pos int
}

/// Helper function.
///
func (t token) debug() {
	fmt.Fprintf(os.Stderr, "%#v\n", t)
}

/// Reads the next token from a scanner. Returns the token.
///
func (s *tokenScanner) scanToken() token {
	for len(s.bytes) > s.pos && s.bytes[s.pos] < 33 {
		s.pos++
	}

	// if at the end, return a comment token
	if len(s.bytes) <= s.pos {
		return token{typ: TOKEN_END, val: ""}
	}

	// get the next character
	c := s.bytes[s.pos]

	// get the next character
	switch {
	case c == ';':
		return s.scanToEnd()
	case c == '[':
		return s.scanEffectiveAddress()
	case c == ',':
		return s.scanOperand()
	case c == '#':
		return s.scanHexLit()
	case c == '%':
		return s.scanBinLit()
	case c == '-':
		return s.scanDecLit()
	case c >= '0' && c <= '9':
		return s.scanDecLit()
	case c >= 'A' && c <= 'Z':
		return s.scanIdentifier()
	case c == '"' || c == '\'' || c == '`':
		return s.scanString(c)
	}

	return s.scanChar()
}

/// Scan a single character.
///
func (s *tokenScanner) scanChar() token {
	i := s.pos

	// advance the scan pos
	s.pos += 1

	// some characters are special tokens
	switch s.bytes[i] {
	case '*':
		return token{typ: TOKEN_HERE}
	}

	return token{typ: TOKEN_CHAR, val: s.bytes[i]}
}

/// Scan to the end of the input and return.
///
func (s *tokenScanner) scanToEnd() token {
	text := string(s.bytes[s.pos:])

	// skip to the end
	s.pos = len(s.bytes)

	// a hard-coded token
	return token{typ: TOKEN_END, val: strings.TrimSpace(text)}
}

/// Scan a comma-separated operand token.
///
func (s *tokenScanner) scanOperand() token {
	s.pos += 1

	// scan the next token as the operand
	t := s.scanToken()

	// make sure there was an operand
	if t.typ == TOKEN_END {
		panic("expected operand")
	}

	return token{typ: TOKEN_OPERAND, val: t}
}

/// Scan a list of comma-separated tokens.
///
func (s *tokenScanner) scanOperands() []token {
	tokens := make([]token, 0, 3)

	// is this the end of the operand list?
	for t := s.scanToken(); t.typ != TOKEN_END; {
		tokens = append(tokens, t)

		// get another token, are we at the end?
		if t = s.scanToken(); t.typ != TOKEN_OPERAND {
			if t.typ == TOKEN_END {
				break
			}

			panic("unexpected token")
		}

		// expand the operand
		t = t.val.(token)
	}

	return tokens
}

/// Scan the effective address of a register (I).
///
func (s *tokenScanner) scanEffectiveAddress() token {
	s.pos += 1

	// scan the next token to take the effective address
	if t := s.scanToken(); t.typ != TOKEN_I {
		panic("illegal indirection")
	}

	// terminate with closing bracket
	if t := s.scanToken(); t.typ != TOKEN_CHAR || t.val.(byte) != ']' {
		panic("illegal effective address")
	}

	return token{typ: TOKEN_EFFECTIVE_ADDRESS}
}

/// Scan an identifier: instruction, register, or label reference.
///
func (s *tokenScanner) scanIdentifier() token {
	i := s.pos

	// advance to the first non-identifier character
	for ;s.pos < len(s.bytes);s.pos++ {
		c := s.bytes[s.pos]

		// validate identifier characters
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			break
		}
	}

	// extract the label
	id := string(s.bytes[i:s.pos])

	// determine whether the label is an instruction, register or reference
	switch id {
	case "V0":
		return token{typ: TOKEN_V, val: 0}
	case "V1":
		return token{typ: TOKEN_V, val: 1}
	case "V2":
		return token{typ: TOKEN_V, val: 2}
	case "V3":
		return token{typ: TOKEN_V, val: 3}
	case "V4":
		return token{typ: TOKEN_V, val: 4}
	case "V5":
		return token{typ: TOKEN_V, val: 5}
	case "V6":
		return token{typ: TOKEN_V, val: 6}
	case "V7":
		return token{typ: TOKEN_V, val: 7}
	case "V8":
		return token{typ: TOKEN_V, val: 8}
	case "V9":
		return token{typ: TOKEN_V, val: 9}
	case "VA":
		return token{typ: TOKEN_V, val: 10}
	case "VB":
		return token{typ: TOKEN_V, val: 11}
	case "VC":
		return token{typ: TOKEN_V, val: 12}
	case "VD":
		return token{typ: TOKEN_V, val: 13}
	case "VE":
		return token{typ: TOKEN_V, val: 14}
	case "VF":
		return token{typ: TOKEN_V, val: 15}
	case "R":
		return token{typ: TOKEN_R}
	case "I":
		return token{typ: TOKEN_I}
	case "F":
		return token{typ: TOKEN_F}
	case "HF":
		return token{typ: TOKEN_HF}
	case "K":
		return token{typ: TOKEN_K}
	case "A":
		return token{typ: TOKEN_ASCII}
	case "D", "DT":
		return token{typ: TOKEN_DT}
	case "S", "ST":
		return token{typ: TOKEN_ST}
	case "CLS", "RET", "EXIT", "LOW", "HIGH", "SCU", "SCD", "SCR", "SCL", "SYS", "JP", "CALL", "SE", "SNE", "SGT", "SLT", "SKP", "SKNP", "LD", "OR", "AND", "XOR", "ADD", "SUB", "SUBN", "MUL", "DIV", "SHR", "SHL", "BCD", "RND", "DRW":
		return token{typ: TOKEN_INSTRUCTION, val: id}
	case "ASCII", "BYTE", "WORD", "ALIGN", "PAD":
		return token{typ: TOKEN_INSTRUCTION, val: id}
	case "BREAK":
		return token{typ: TOKEN_BREAK}
	case "ASSERT":
		return token{typ: TOKEN_ASSERT}
	case "EQU":
		return token{typ: TOKEN_EQU}
	case "VAR":
		return token{typ: TOKEN_VAR}
	case "SUPER":
		return token{typ: TOKEN_SUPER}
	case "EXTENDED":
		return token{typ: TOKEN_EXTENDED}
	}

	if i == 0 {
		return token{typ: TOKEN_LABEL, val: id}
	} else {
		return token{typ: TOKEN_ID, val: id}
	}
}

/// Scan a decimal literal.
///
func (s *tokenScanner) scanDecLit() token {
	i := s.pos

	// skip a unary minus negation
	if s.bytes[i] == '-' {
		s.pos += 1
	}

	// find the first non-numeric character
	for ;s.pos < len(s.bytes);s.pos += 1 {
		if strings.IndexByte("0123456789", s.bytes[s.pos]) < 0 {
			break
		}
	}

	// convert the hex value to a signed number
	if n, err := strconv.ParseInt(string(s.bytes[i:s.pos]), 10, 32); err == nil {
		return token{typ: TOKEN_LIT, val: int(n)}
	}

	panic(fmt.Errorf("illegal decimal value: %s", string(s.bytes[i:s.pos])))
}

/// Scan a hexadecimal literal.
///
func (s *tokenScanner) scanHexLit() token {
	i := s.pos

	// find the first non-hex character
	for s.pos += 1;s.pos < len(s.bytes);s.pos += 1 {
		if strings.IndexByte("0123456789ABCDEF", s.bytes[s.pos]) < 0 {
			break
		}
	}

	// convert the hex value to an unsigned number
	if n, err := strconv.ParseInt(string(s.bytes[i+1:s.pos]), 16, 32); err == nil {
		return token{typ: TOKEN_LIT, val: int(n)}
	}

	panic(fmt.Errorf("illegal hex value: #%s", string(s.bytes[i:s.pos])))
}

/// Scan a binary literal.
///
func (s *tokenScanner) scanBinLit() token {
	i := s.pos

	// find the first non-binary character
	for s.pos += 1;s.pos < len(s.bytes);s.pos += 1 {
		if strings.IndexByte(".01", s.bytes[s.pos]) < 0 {
			break
		}
	}

	// replace all '.' with '0'
	v := strings.Replace(string(s.bytes[i+1:s.pos]), ".", "0", -1)

	// convert the hex value to an unsigned number
	if n, err := strconv.ParseInt(v, 2, 32); err == nil {
		return token{typ: TOKEN_LIT, val: int(n)}
	}

	panic(fmt.Errorf("illegal binary value: $%s", string(s.bytes[i:s.pos])))
}

/// Scan a quoted string.
///
func (s *tokenScanner) scanString(term byte) token {
	s.pos += 1

	// store starting position
	i := s.pos

	// find the terminating quotation
	for s.pos < len(s.bytes) && s.bytes[s.pos] != term {
		s.pos++
	}

	// advance past the terminator
	s.pos += 1

	return token{typ: TOKEN_TEXT, val: string(s.bytes[i:s.pos-1])}
}
