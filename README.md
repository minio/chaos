# Chaos
A framework for testing Minio's fault tolerance capability.

# Design

 - Contains master and worker programs. Worker has be run on every node along with the Minio server and Master commands the chaos-workers on the remote nodes.
 - The master program pings the chaos workers across nodes and ensures that they are running and the choas workers in turn verify that the Minio server is running on the remote node at the specified port during initialization.
 - On the event of any chaos worker not reachable or if the Minio process is not running on the remote node the test fails.
 - Currrently RoundRobin Chaos Test is introduced wherein Minio server is stopped on each node and recovered after the specified recovery time one after the other.
   
# How to run. 

- Fetch the project 

```sh
$ go get -d github.com/minio/chaos` 
```

- Build master program.

```sh
$ cd $GOPATH/src/github.com/minio/chaos/master && go build
```

- Run Minio Distributed server using systemd script on remote nodes. The workers use `systemd` to control the Minio process. [Click here](https://github.com/minio/minio/tree/master/dist/linux-systemd/distributed) for info on configuring systemd to run Minio Distributed.
- Build and run chaos-workers on each these remote nodes. Use `sudo` to run worker. Need privilaged access to control Minio process using systemd.
  
```sh
$ go get github.com/minio/chaos`
$ cd $GOPATH/src/github.com/minio/chaos/worker/ && go build`
$ sudo ./worker
```

- Run master. 

```sh
$ master -endpoints="<Node-1-IP>:9997,<NODE-2-IP>:9997,<NODE-3-IP>:9997...... -recover=30 -rounds=10 -minio-port=9199"
```

Chaos workers will be running at port "9997" now.   
    

# Options

```sh
-recover : Removery time after the failure injection on remote Minio node.
```

```sh
-endpoints: "," separted <IP>:<PORT> at which Remote Chaos Workers are Running".
```

```sh
-rounds: Number of rounds the chaos test has to be run.
```

```sh
-minio-port: Port at which Minio server is running on remote node. Port 9000 is taken as the default value if no value is provided. 
```

# TODO

- By Extending the `Chaos` interface different failures can be introduced. Next step is to make it more easier to run, with more options and different failures like multiple node failures will introduced.
