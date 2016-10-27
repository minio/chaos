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
	"encoding/json"
	"net/http"
	"sync"
	"time"
	//	"github.com/prometheus/client_golang/prometheus"
)

const (
	serverUp   = "Server Running"
	serverDown = "Server Stopped by Worker"
)

// experimenting with prometheus for counters in the project.
// will move onto use other features it this looks good.1
//var StoppedNodesCounter = prometheus.NewCounter(
//	prometheus.CounterOpts{
//		Namespace: "Minio",
//		Subsystem: "Minio-Chaos-Test",
//		Name:      "StoppedNodesCounter",
//		Help:      "Total number of Nodes Stopped by Chaos workers currently",
//	})

// register the prometheus service.
func init() {
	//prometheus.MustRegister(StoppedNodesCounter)
}

// Report - contains fields to keep the status of the chaos test across all chaos workers on all nodes.
// Here are some of the expected fields in the report.
// Since - Time since the test is running.
// ServerStatus - Contains the status of the Minio server on the nodes where the chaos worker is running.
// 		  Since the chaos worker stops the server and restarts after a definate interval, the status helps to whether
//                the Minio server is running or has been stopped by the
// 		  chaos worker.
type Report struct {
	// Lock to be used before mutating any field of the status.
	sync.Mutex
	// Add all necessary fields to report the status of the chaos operations.
	TestRunningSince time.Time
	// Status of all the Minio servers.
	// Obtaining by visiting `/report` of the master process of the chaos test.
	Status []ServerStatus
}

// ServerStatus - Field to store the status of the Minio server running on the nodes.
type ServerStatus struct {
	// Addr of the Minio Server Node.
	MinioAddr string
	// status is either "serverUp" or "serverDown".
	Status string
	// Time since the above status holds.
	StatusSince time.Time
	// The next expected change in status of the Node.
	NextChange string
	// Time aftedr which the status change can be expected.
	NextChangeIn time.Time
}

// Wraps the status report of the test and exposes HTTP handler interface for
// viewing the status.

type reportHandler struct {
	report *Report
}

// HTTP handler which fetches the status of all the nodes during the chaos test.
// Pings all the chaos workers and obtain the status of the Minio nodes.
// Once the status obtained write it to as a HTTP response.
func (rep reportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	en := json.NewEncoder(w)

	rep.report.Lock()
	defer rep.report.Unlock()

	// TODO: Ping all the chaos workers and obtain the status of the Minio nodes.
	// Once the status obtained write it to as a HTTP response.

	if err := en.Encode(Report{

	// other parameters to be reported goes here.
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
