// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package context

import (
	"github.com/m3db/m3x/pool"
	"github.com/m3db/m3x/resource"
)

// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// finalizersArrayPool provides a pool for resourceFinalizer slices.
type finalizersArrayPool interface {
	// Init initializes the array pool, it needs to be called
	// before Get/Put use.
	Init()

	// Get returns the a slice from the pool.
	Get() []resource.Finalizer

	// Put returns the provided slice to the pool.
	Put(elems []resource.Finalizer)
}

type finalizersFinalizeFn func([]resource.Finalizer) []resource.Finalizer

type finalizersArrayPoolOpts struct {
	Options     pool.ObjectPoolOptions
	Capacity    int
	MaxCapacity int
	FinalizeFn  finalizersFinalizeFn
}

type finalizersArrPool struct {
	opts finalizersArrayPoolOpts
	pool pool.ObjectPool
}

func newFinalizersArrayPool(opts finalizersArrayPoolOpts) finalizersArrayPool {
	if opts.FinalizeFn == nil {
		opts.FinalizeFn = defaultFinalizersFinalizerFn
	}
	p := pool.NewObjectPool(opts.Options)
	return &finalizersArrPool{opts, p}
}

func (p *finalizersArrPool) Init() {
	p.pool.Init(func() interface{} {
		return make([]resource.Finalizer, 0, p.opts.Capacity)
	})
}

func (p *finalizersArrPool) Get() []resource.Finalizer {
	return p.pool.Get().([]resource.Finalizer)
}

func (p *finalizersArrPool) Put(arr []resource.Finalizer) {
	arr = p.opts.FinalizeFn(arr)
	if max := p.opts.MaxCapacity; max > 0 && cap(arr) > max {
		return
	}
	p.pool.Put(arr)
}

func defaultFinalizersFinalizerFn(elems []resource.Finalizer) []resource.Finalizer {
	var empty resource.Finalizer
	for i := range elems {
		elems[i] = empty
	}
	elems = elems[:0]
	return elems
}

type finalizersArr []resource.Finalizer

func (elems finalizersArr) grow(n int) []resource.Finalizer {
	if cap(elems) < n {
		elems = make([]resource.Finalizer, n)
	}
	elems = elems[:n]
	// following compiler optimized memcpy impl
	// https://github.com/golang/go/wiki/CompilerOptimizations#optimized-memclr
	var empty resource.Finalizer
	for i := range elems {
		elems[i] = empty
	}
	return elems
}
