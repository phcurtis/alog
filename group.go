// Copyright 2017 phcurtis grplog Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grplog

import (
	"fmt"
	"io"
)

// SetFlags - sets stdlib log flags for all group log levels to fval.
func (g *GlvlStruct) SetFlags(fval int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.lvlList() {
		(*v.level).log.SetFlags(fval)
	}
}

// SetIgnore - set each level individual ignore state.
func (g *GlvlStruct) SetIgnore(state bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.lvlList() {
		(*v.level).ignore = state
	}
}

// GetIgnoreAll - return ignoreall flag for group.
func (g *GlvlStruct) GetIgnoreAll() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.ignoreall
}

// SetIgnoreAll - sets ignoreall flag for group.
func (g *GlvlStruct) SetIgnoreAll(state bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.ignoreall = state
}

// SetLabel - applies group string to all grplog levels.
func (g *GlvlStruct) SetLabel(glabel string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.lvlList() {
		(*v.level).log.SetPrefix(glabel + v.Blab)
	}
}

// SetPkgFlags - set group
func (g *GlvlStruct) SetPkgFlags(f int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.lvlList() {
		(*v.level).flags = f
	}
}

// SetOutput .. sets io.Writer for each log level in group.
func (g *GlvlStruct) SetOutput(w io.Writer) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.lvlList() {
		(*v.level).log.SetOutput(w)
		(*v.level).logOutput = w
	}
}

// Println - calls Println for each level of the group with args passed in.
func (g *GlvlStruct) Println(x ...interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, v := range g.lvlList() {
		if (*v.level).anyIgnore(false) {
			continue
		}
		_ = (*v.level).outll(0, fmt.Sprintln(x...))
	}
}

// CondPrintln - conditional version of Println
func (g *GlvlStruct) CondPrintln(cond bool, x ...interface{}) {
	if cond {
		g.mu.Lock()
		defer g.mu.Unlock()
		for _, v := range g.lvlList() {
			if (*v.level).anyIgnore(false) {
				continue
			}
			_ = (*v.level).outll(0, fmt.Sprintln(x...))
		}
	}
}
