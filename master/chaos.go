/*
 * Minio Cloud Storage, (C) 2015, 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"time"
)

// Interface to which any choas function should satisfy.
type Chaos interface {
	// Responsible for calling the relevant failure method on the remote node.
	// Should result in failure injection on the Minio server on the remote node.
	Fail(c ChaosWorker) error
	// Responsible for recovering the previsouly injected failure on the remote node.
	Recover(c ChaosWorker) error
	// Executes
	Execute(t ChaosTest) error
}

// Has Methods,
// * `Fail` - To stop Minio server on the remote node.
// * `Recover` - To start Minio server on the remote node.
// ^^ These are the most basic Fail and Recover functions.
type GenericFail struct {
}

func (fail GenericFail) Fail(c ChaosWorker) error {
	var err error
	// Address containing info of the port at which Minio has to be run
	// on the remote node of the worker.
	minioRemoteAddr := c.Node.Addr
	args := &minioRemoteAddr
	// expecting only error response from the RPC call.
	// not relying on the reply.
	reply := struct{}{}
	// Stop the Minio server on the remote node.
	log.Println("Attempting to stop Minio server on node: ", c.WorkerEndpoint)
	err = c.Client.Call("ChaosWorker.StopMinioServer", args, &reply)
	if err == nil {
		log.Println("Minio server stopped on node: ", c.WorkerEndpoint)
	}
	return err
}

// Most basic recovery on the remote node,
// Starts the Minio server on the remote node.
// Used to used to recover after the `Fail` method is called.
func (fail GenericFail) Recover(c ChaosWorker) error {
	var err error
	// Address containing info of the port at which Minio has to be run
	// on the remote node of the worker.
	minioRemoteAddr := c.Node.Addr
	args := &minioRemoteAddr
	// expecting only error response from the RPC call.
	// not relying on the reply.
	reply := struct{}{}

	log.Println("Attempting to Start Minio server on node: ", c.WorkerEndpoint)
	// start the Minio server on the remote node.
	err = c.Client.Call("ChaosWorker.StartMinioServer", args, &reply)
	// Successfully started Minio server on the remote node,
	// log the result.
	if err == nil {
		log.Println("Minio server Started on node: ", c.WorkerEndpoint)
	}
	return err
}

// Introduce Failure and Remove each remote node one by one.
// There's a delay of specified time period between failure and removery of a node.
type RoundRobinChaos struct {
	GenericFail
}

// Iterates over all Choas workers,
func (round *RoundRobinChaos) Execute(chaos ChaosTest) error {
	var err error
	for _, worker := range chaos.ChaosWorkers {
		// TODO: Remove Sleep and introduce Context/Timeout and signalling methods to improve the logic.
		// Call the Fail method to introduce failure on the remote node.
		err = round.Fail(*worker)
		if err != nil {
			return err
		}
		// TODO: Use a better logic than just to sleep.
		// Sleep for the time interval set for the failure recovery.
		time.Sleep(time.Duration(chaos.RecoveryTime) * time.Second)
		// Call the `Recover` Method to recover from the previously introduced failure on the remote node.
		err = round.Recover(*worker)
		if err != nil {
			return err
		}
	}
	return nil
}
