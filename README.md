# Chaos
A framework for testing Minio's fault tolerance capability.

# Design

 - With the help of docker-compose (see docker-compose.yml), Minio distributed with 4 containers started. Each of these containers also has the worker program running.
 - Minio servers are exposed on host port 9010-9013, and the workers on 9014-9017.
 - Workers use `supervisord` for controlling minio process in each of the containers.
 - The master program pings the chaos workers across nodes and ensures that they are running and the chaos workers in turn verify that the Minio server is running on the remote node at the specified port during initialization.
 - On the event of any chaos worker not reachable or if the Minio process is not running on the remote node the test fails.
 - Currrently RoundRobin Chaos Test is introduced wherein Minio server is stopped on each node and recovered after the specified recovery time one after the other.
 - The load has to generated externally against the Minio server, the chaos project just provides easy mechanism to start, control and monitor Minio process.
 - The Minio server and supervisord logs for Minio and the worker process will be persisted in the logs folder. 
 - As the Minio server run during choas setup is meant for fault tolerant and load testing the backend from the containers are not persisted. This means its not necessary to perform the cleanups after your tests. Make sure any check on the backend is done by logging into the container while the compose is still running. See the use case section at the end of the document.
   
# How to run. 

- Fetch the project 

  ```sh
  $ git clone https://github.com/minio/chaos.git && cd chaos
  ```

- Run Minio distributed and worker process using docker-compose. 
  - Minio servers will be exposed on ports 9010-9013 for any load testing.
  - The worker RPC servers will be available on ports 9014-9017 on the host.

  ```sh
  $ docker-compose up
  ```

- Build master program.

  ```sh
  $ cd master && go build

  ```

- Run master. Specify address of only those nodes on which you want to `stop->wait->start` Minio servers. 

  ```sh
  $ master -endpoints="<Node-1-IP>:9997,<NODE-2-IP>:9997,<NODE-3-IP>:9997......" -recover=30 -rounds=10 
  ```

# Options for master.

```sh
-endpoints: "," separted <IP>:<PORT> at which Remote Chaos Workers are Running".
```

```sh
-recover : Removery time after the failure injection on remote Minio node.
```

```sh
-rounds: Number of rounds the chaos test has to be run.
```

# Controlling specific nodes.
- Using master program.

  Master program can be control specific nodes. The following example stops Minio server node 4 for 200 seconds.
  Get the default setup running using `docker-compose up` before the master command is executed.

  ```sh
  $ master -endpoints="127.0.0.1:9017" -recover=200 -rounds=1
  ```

- Manually control.

  You can also directly use `supervisorctl` from within the container to control Minio process. 
  The following example permanently stops the Minio server on node 1.
  Get the default setup running using `docker-compose up` before the master command is executed.

  ```sh
  docker-compose exec minio1 supervisorctl stop minio

  ```

  Start Minio node again by using 
  ```sh
  docker-compose exec minio1 supervisorctl start minio
  ```

# Example use case.

Using chaos for `mc admin heal` testing.

- Run the setup.
  ```sh
  $ docker-compose up
  ```

- Stop one of the node.
  ```sh
  $ docker-compose exec minio4 supervisorctl stop minio
  ```

- Uplaod a file from host.
  ```sh
  $ mc config host add minio4 http://127.0.0.1:9013 minio123 minio12345
  $ mc mb minio4/test-bucket/
  $ mc cp 1gb minio3/test-bucket/
  ```

- Explore the backend (the backend is available at `/export` within the container) and try `mc admin heal`, `mc` is installed and configured on each of the container.
  ```sh
   $ docker-compose exec minio3 /bin/bash
   # ls /export/ 
   # mc admin heal list myminio/test-bucket
  ```

- Choas makes it easy to start, control, manage, monitor fault tolerant related tests for Minio.

# TODO

- By Extending the `Chaos` interface different failures can be introduced. Next step is to make it more easier to run, with more options and different failures like multiple node failures will introduced.
