/* 
 * Copyright (c) 2018, Giovanni Dante Grazioli <deroad@libero.it>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package main

type R2Channel chan *R2Test
type R2Results chan *TestResult

type R2Pool struct {
	Size int
	Tests R2Channel
	Results R2Results
}

func R2Routine(pool *R2Pool, done chan bool) {
	for {
		select {
		case test := <- pool.Tests:
			pool.Results <- test.Exec()
		default:
			done <- true
			return
		}
	}
}

func (pool R2Pool) PerformTests(ptests *[]R2Test) bool {
	success := true
	tests := (*ptests)
	done := make(chan bool, pool.Size)
	length := len(tests)
	pool.Tests = make(R2Channel, length)
	pool.Results = make(R2Results, length)
	for index := range tests {
		pool.Tests <- &tests[index]
	}

	for i := 0; i < pool.Size; i++ {
		go R2Routine(&pool, done)
	}
	for i := 0; i < pool.Size; i++ {
		<- done
	}
	for i := 0; i < length; i++ {
		result := <- pool.Results
		if result.Print(true) {
			success = false
		}
	}
	return success
}

func NewR2Pool(size int) *R2Pool {
	return &R2Pool{size, nil, nil}
}
