// Package primative provides functions to be able to build primative recursive functions.
package primative

import (
	"fmt"
	"strings"
)

func format(f Func, s fmt.State, verb rune) {
	type verboser interface{
		Verbose(prec int) string
	}

	switch verb {
	case 'v':
		if prec, ok := s.Precision(); ok {
			if prec == 0 {
				fmt.Fprint(s, f.(fmt.Stringer).String())
				return
			}

			if f, ok := f.(verboser); ok {
				fmt.Fprintf(s, f.Verbose(prec))
				return
			}
		}

		if s.Flag('+') {
			if f, ok := f.(verboser); ok {
				fmt.Fprint(s, f.Verbose(-1))
				return
			}
		}

		fallthrough

	case 's':
		fmt.Fprint(s, f.(fmt.Stringer).String())
	}
}

func argString(x []uint) string {
	var list []string

	for i := range x {
		list = append(list, fmt.Sprint(x[i]))
	}

	return strings.Join(list, ", ")
}

type named struct {
	name string
	Func
}

func Name(name string, f Func) Func {
	return &named{
		name: name,
		Func: f,
	}
}

func (f *named) String() string {
	return f.name
}

func (f *named) Verbose(prec int) string {
	if prec == 0 {
		return f.String()
	}

	if prec < 0 {
		return fmt.Sprintf("%+v", f.Func)
	}

	return fmt.Sprintf("%.*v",  prec-1, f.Func)
}

func (f *named) Format(s fmt.State, verb rune) {
	format(f, s, verb)
}

func (f *named) Compose(g ...Func) Func {
	return Compose(f, g...)
}

func Debug(f Func) Func {
	return &funcN{
		ary: f.Ary(),
		f: func(x ...uint) uint {
			args := argString(x)

			result := f.Apply(x...)

			fmt.Printf("%v(%s) = %d\n", f, args, result)

			return result
		},
		s: fmt.Sprint(f),
	}
}

func Panic(v interface{}) Func {
	s := fmt.Sprint(v)

	return &funcN{
		ary: 0,
		f: func(x ...uint) uint {
			panic(s)
		},
		s: fmt.Sprintf("panic(%q)", s),
	}
}
