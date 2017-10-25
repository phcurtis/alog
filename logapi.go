// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Core functionality of these included func APIs were derived from stdlib log (as of go 1.9)
// an as such We are including the enclosed contents of this file to be part of that
// work as such with the following credits: to Copyright 2009 The Go Authors. All rights reserved.

package grplog

import (
	"fmt"
)

// Fatal - stdlib.log level Fatal but possible decorate, etc.
func (l *lvlStruct) Fatal(x ...interface{}) {
	l.outExit(fmt.Sprint(x...))
}

// Panic - stdlib.log level Panic but possible decorate, etc.
func (l *lvlStruct) Panic(x ...interface{}) {
	l.outPanic(fmt.Sprint(x...))
}

// Print - stdlib.log level Print but possible ignore or decorate, etc.
func (l *lvlStruct) Print(x ...interface{}) {
	if l.anyIgnore(true) {
		return
	}
	_ = l.out(fmt.Sprint(x...))
}

// Fatalf stdlib.log Fatalf but possible decorate, etc.
func (l *lvlStruct) Fatalf(f string, x ...interface{}) {
	l.outExit(fmt.Sprintf(f, x...))
}

// Panicf - stdlib.log level Panicf but possible decorate, etc.
func (l *lvlStruct) Panicf(f string, x ...interface{}) {
	l.outPanic(fmt.Sprintf(f, x...))
}

// Printf - stdlib.log level Printf but possible ignore or decorate, etc.
func (l *lvlStruct) Printf(f string, x ...interface{}) {
	if l.anyIgnore(true) {
		return
	}
	_ = l.out(fmt.Sprintf(f, x...))
}

// Fatalln - stdlib.log level Fatalln but possible decorate, etc.
func (l *lvlStruct) Fatalln(x ...interface{}) {
	l.outExit(fmt.Sprintln(x...))
}

// Panicln - stdlib.log level "Panicln" but possible decorate, etc.
func (l *lvlStruct) Panicln(x ...interface{}) {
	l.outPanic(fmt.Sprintln(x...))
}

// Println - stdlib.log level "Println" but possible ignore or decorate, etc.
func (l *lvlStruct) Println(x ...interface{}) {
	if l.anyIgnore(true) {
		return
	}
	_ = l.out(fmt.Sprintln(x...))
}
