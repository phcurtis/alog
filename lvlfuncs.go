// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog

import (
	"fmt"
	"io"
)

// AnyIgnore - returns true if level or parent group has "ignore" set
func (l *LvlStruct) anyIgnore(protect bool) bool {
	if protect {
		l.mu.Lock()
		defer l.mu.Unlock()
	}
	return l.ignore || (l.par != nil && l.par.ignoreall)
}

// CondPrint - conditional version of Print
func (l *LvlStruct) CondPrint(cond bool, x ...interface{}) {
	if cond {
		if l.anyIgnore(true) {
			return
		}
		_ = l.out(fmt.Sprint(x...))
	}
}

// CondPrintf ... conditional version of Printf
func (l *LvlStruct) CondPrintf(cond bool, f string, x ...interface{}) {
	if cond {
		if l.anyIgnore(true) {
			return
		}
		_ = l.out(fmt.Sprintf(f, x...))
	}
}

// CondPrintln - conditional version of Println
func (l *LvlStruct) CondPrintln(cond bool, x ...interface{}) {
	if cond {
		if l.anyIgnore(true) {
			return
		}
		_ = l.out(fmt.Sprintln(x...))
	}
}

// AlignFile - return alignment [minimum width] for filename stuff
func (l *LvlStruct) AlignFile() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.align.filea
}

// SetAlignFile - set alignment [minimum width] for filename stuff
func (l *LvlStruct) SetAlignFile(minWidth int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if minWidth > LogAlignFileMax {
		minWidth = LogAlignFileMax
	} else if minWidth < 0 {
		minWidth = 0
	}
	l.align.filea = minWidth
}

// AlignFunc - return alignment [minimum width] for funcname stuff
func (l *LvlStruct) AlignFunc() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.align.funca
}

// SetAlignFunc - set alignment [minimum width] for funcname stuff
func (l *LvlStruct) SetAlignFunc(minWidth int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if minWidth > LogAlignFuncMax {
		minWidth = LogAlignFuncMax
	} else if minWidth < 0 {
		minWidth = 0
	}
	l.align.funca = minWidth
}

// Flags - returns the log flags.
func (l *LvlStruct) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.log.Flags()
}

// SetFlags - sets the log flags.
func (l *LvlStruct) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.SetFlags(flag)
}

// SetIgnore - set log ignore state.
func (l *LvlStruct) SetIgnore(b bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ignore = b
}

// Ignore - returns log ignore state.
func (l *LvlStruct) Ignore() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.ignore
}

// GetOutput - returns the log Output io.Writer
func (l *LvlStruct) GetOutput() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logOutput
}

// SetOutput ... calls stdlib log.SetOutput.
// For group level best to configure as needed during creation
// see NewSpecial func.
func (l *LvlStruct) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.SetOutput(w)
	l.logOutput = w
}

// Prefix - returns 'prefix' label.
func (l *LvlStruct) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.log.Prefix()
}

// SetPrefix - set prefix for log level.
// (group) SetLabel if called will revert all base labels to
// Compile time base labels. This may change in the future.
// It was also thought of removing this function however
// one may want to modify the Gtrace prefix since it doesn't
// belong to a Group.
func (l *LvlStruct) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.SetPrefix(prefix)
}

// PkgFlags -
func (l *LvlStruct) PkgFlags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flags
}

// SetPkgFlags -
func (l *LvlStruct) SetPkgFlags(f int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flags = f
}
