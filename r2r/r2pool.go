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

type TestsOptions struct {
	Debug      bool
	Sequence   bool
	ErrorsOnly bool
	Jobs       int
}

type R2Pool struct {
	Tests   R2Channel
	Results R2Results
	Options *TestsOptions
}

func R2Routine(pool *R2Pool, done chan bool) {
	for {
		select {
		case test := <-pool.Tests:
			pool.Options.Println("Executing", test.Name)
			pool.Results <- test.Exec(pool.Options)
			pool.Options.Println("Result returned.")
		default:
			done <- true
			return
		}
	}
}

func (pool R2Pool) PerformTests(regressions *R2RegressionTest) bool {
	success := true
	tests := regressions.Tests
	done := make(chan bool, pool.Options.Jobs)
	length := len(tests)
	pool.Tests = make(R2Channel, length)
	pool.Results = make(R2Results, length)

	pool.Options.Println("Preparing", length, "tests..")

	for index := range tests {
		pool.Tests <- &tests[index]
	}

	if pool.Options.Jobs > 1 {
		pool.Options.Println("Poolsize:", pool.Options.Jobs)
		for i := 0; i < pool.Options.Jobs; i++ {
			go R2Routine(&pool, done)
			pool.Options.Println("Waiting end of tests...")
			for i := 0; i < pool.Options.Jobs; i++ {
				<-done
			}
		}
	} else {
		pool.Options.Println("single-thread mode")
		R2Routine(&pool, done)
	}

	for i := 0; i < length; i++ {
		result := <-pool.Results
		if !pool.Options.Sequence {
			result.Print(true)
		}
		if !result.Success && !result.Test.Broken {
			success = false
		}
	}
	return success
}

func NewR2Pool(options *TestsOptions) *R2Pool {
	return &R2Pool{nil, nil, options}
}
