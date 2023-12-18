package vdf

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type KeyMap struct {
	k string
	m map[string]interface{}
}

func ParseUtf8(reader io.Reader) map[string]interface{} {

	br := bufio.NewReader(reader)

	m := make(map[string]interface{})
	km := []KeyMap{{"", m}}
	var l int
	var k = true

	var r rune
	var err error
	for {
		r, _, err = br.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		if r == '"' {
			l++
			if k {
				km = append(km, KeyMap{})
			}
			var c string
			for {
				r, _, _ = br.ReadRune()
				cl := len(c)
				if r == '"' && (cl == 0 || c[cl-1] != '\\') {
					break
				}
				c += string(r)
			}
			if k {
				km[l].k = c
			} else {
				km[l-1].m[km[l].k] = c
				km = km[:l]
			}
			l--
			k = !k
			continue
		}

		if r == '{' {
			l++
			if sk, ok := km[l-1].m[km[l].k]; ok {
				km[l].m = sk.(map[string]interface{})
			} else {
				km[l].m = make(map[string]interface{})
			}
			k = true
			continue
		}

		if r == '}' {
			km[l-1].m[km[l].k] = km[l].m
			km = km[:l]
			l--
		}
	}
	return m
}

func ParseUtf16(reader io.Reader) map[string]interface{} {

	br := bufio.NewReader(reader)

	m := make(map[string]interface{})
	km := []KeyMap{{"", m}}
	var l int
	var k = true

	var r rune
	var r2 rune
	var err error
	for {
		r, _, _ = br.ReadRune()
		r2, _, err = br.ReadRune()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		if r == '"' && r2 == 0 {
			l++
			if k {
				km = append(km, KeyMap{})
			}
			var c string
			for {
				r, _, _ = br.ReadRune()
				r2, _, _ = br.ReadRune()
				cl := len(c)
				if r == '"' && r2 == 0 && (cl == 0 || c[cl-1] != '\\') {
					break
				}
				c += string(r)
			}
			if k {
				km[l].k = strings.ToLower(c)
			} else {
				km[l-1].m[km[l].k] = c
				km = km[:l]
			}
			l--
			k = !k
			continue
		}

		if r == '{' && r2 == 0 {
			l++
			if sk, ok := km[l-1].m[km[l].k]; ok {
				km[l].m = sk.(map[string]interface{})
			} else {
				km[l].m = make(map[string]interface{})
			}
			k = true
			continue
		}

		if r == '}' && r2 == 0 {
			km[l-1].m[km[l].k] = km[l].m
			km = km[:l]
			l--
		}

		if r == '/' && r2 == 0 {
			for {
				r, _, _ = br.ReadRune()
				r2, _, _ = br.ReadRune()
				if r == '\n' && r2 == 0 {
					break
				}
			}
		}
	}

	return m
}
