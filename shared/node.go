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

// Package  containing shared data structure between master and the workers.
package shared

// MinioNode - info of the Minio node under chaos test.
type MinioNode struct {
	// flag indicating whether Minio node has to started.
	// unless specifed with a `-start` flag an attempt to start Minio server
	// on the remote node is done.`
	StartMinio bool
	// Adress of format 127.0.0.1:<PORT>,
	// Used to validate whether Minio server is running on each of the nodes
	// in the specified port.
	Addr string
}
