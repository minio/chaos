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
	"flag"
	"log"
	"strconv"
	"strings"
)

// verify whether the Minio server port is a valid integer.
func verifyMinioPort(portStr string) error {
	_, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// TODO: parse all the flags here.
	endPointStr := flag.String("endpoints", "", "RPC endpoints of workers.")
	recoverStr := flag.String("recover", "10", "Recovery time of the remote node after Choas.")
	roundsStr := flag.String("rounds", "1", "Number of rounds the Choas test has to be run.")
	// Port at which Minio server is run on remote nodes.
	// Default Minio server port of 9000 is taken as the default option.
	minioPortStr := flag.String("minio-port", "9000", "Port at which Minio server is run on remote nodes.")
	// parse the command line flags.
	flag.Parse()
	// obtain all the node end points.
	endPoints := strings.Split(*endPointStr, ",")
	// verify whether a valid integer port is given.
	if err := verifyMinioPort(*minioPortStr); err != nil {
		log.Fatalf("Please enter a valid integer port at which Minio server is running on remote nodes.")
	}

	// minioAddr - will be used by workers on the respective remote nodes to verify whether Minio is running as a service
	// before the chaos test is started.
	minioAddr := "http://127.0.0.1:" + *minioPortStr
	// Allocate memory for ChaosWorkers.
	// Allocating one for each node.
	chaosWorkers := make([]*ChaosWorker, len(endPoints))

	// Iterate through the endPoints and create `ChaosTest` instance.
	for i, endPoint := range endPoints {
		worker := ChaosWorker{
			WorkerEndpoint: endPoint,
			Node: MinioNode{
				Addr: minioAddr,
			},
			//TODO: Make use of report Dir.
			ReportDir: "/not-used-yet",
		}
		// push all the workers into the array.
		chaosWorkers[i] = &worker
	}

	// parse the interger value of the recovery string.
	// log and exit in case of an invalid value.
	recoveryTime, err := strconv.Atoi(*recoverStr)
	if err != nil {
		log.Fatalf("Please enter valid time string for recovery: ", err)
	}

	// parse the interger value of the rounds string.
	// log and exit in case of an invalid value.
	rounds, err := strconv.Atoi(*roundsStr)
	if err != nil {
		log.Fatalf("Please enter valid time string for recovery: ", err)
	}
	// Create `ChaosTest` instance here.
	chaosTest := ChaosTest{
		ChaosWorkers: chaosWorkers,
		RecoveryTime: recoveryTime,
	}

	// Initialize all the workers on remote nodes.
	// also confirms that minio server instances are running on the remote nodes.
	if isFailed := chaosTest.InitChaosTest(); isFailed {
		log.Fatal("Iniitalizing of Chaos test failed.")
	}

	log.Println("Initialization finished, Starting Chaos test.")

	// For extending the tests, any new chaos test has to satisfy `Chaos` interface,
	// `RoundRobinChaos` satisfies the `Chaos` interface and it Fails the nodes and recovers them
	// one after another in round robin fashion.
	roundRobinChaos := &RoundRobinChaos{}
	// Run the test for specified number of rounds.
	for j := 0; j < rounds; j++ {
		log.Println("Starting Chaos test... Round ", j+1)
		// Unleash the chaos test.
		err = chaosTest.UnleashChaos(roundRobinChaos)
		// log and exit in case of error.
		if err != nil {
			log.Fatal("Chaos test failed with error: ", err)
		}
	}
}
