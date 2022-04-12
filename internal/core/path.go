package core

import (
	"fmt"
	"strconv"
)

func ProcessPath(path string) (elements []PathElement, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	elements = make([]PathElement, 0)

	ignoreNextSpecial := false
	inBrackets := false
	bracketStart := 0
	inQuotes := false
	quoteStart := 0

	pathRunes := []rune(path)
	pathLen := len(pathRunes)
	elemBuf := make([]rune, 0)

	collect := func(t func([]rune) PathElement) {
		elem := t(elemBuf)
		elements = append(elements, elem)
		elemBuf = make([]rune, 0)
	}

	obj := func(r []rune) PathElement { return NewObjectPathElement(string(r)) }
	arr := func(r []rune) PathElement {
		v, err := strconv.ParseUint(string(r), 10, 64)
		if err != nil {
			panic(err)
		}
		return NewArrayPathElement(uint(v))
	}

	for idx, char := range pathRunes {
		if ignoreNextSpecial {
			// just add char to the buffer and continue
			elemBuf = append(elemBuf, char)
			ignoreNextSpecial = false
			continue
		}

		if char == '\\' {
			ignoreNextSpecial = true
			continue
		}

		if char == '[' && !inQuotes {
			if inBrackets {
				err = fmt.Errorf("unclosed bracket starting at idx %d", bracketStart)
				return
			}

			// Collect object path if we have anything in the buf
			if len(elemBuf) > 0 {
				collect(obj)
			}

			inBrackets = true
			bracketStart = idx
			continue
		}

		if char == ']' && !inQuotes {
			if !inBrackets {
				err = fmt.Errorf("unexpected bracket close at idx %d", idx)
				return
			}
			inBrackets = false

			if len(elemBuf) < 1 {
				err = fmt.Errorf("unexpected bracket close at idx %d", idx)
				return
			}
			collect(arr)

			continue
		}

		if char == '"' && !inBrackets {
			closed := inQuotes
			inQuotes = !inQuotes

			if closed {
				collect(obj)
			} else {
				quoteStart = idx + 1
			}

			continue
		}

		if char == '.' && !inQuotes && !inBrackets {
			if len(elemBuf) > 0 {
				collect(obj)
			}

			continue
		}

		elemBuf = append(elemBuf, char)
		if idx == pathLen-1 {
			collect(obj)
		}
	}

	// If we're still inside quote region...
	if inQuotes {
		// it's time to start yelling
		err = fmt.Errorf("unclosed quote starting at idx %d", quoteStart)
		return
	}

	if inBrackets {
		err = fmt.Errorf("unclosed bracket starting at idx %d", bracketStart)
		return
	}

	return
}
