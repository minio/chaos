# minio-chaos
Chaos framework for testing Minio's fault tolerance capability.


# Initial Design 

 - Contains master and worker programs. Worker has be run on every node along with the Minio server and 
   Master commands the chaos-workers on the remote nodes.
 - The master program pings the chaos workers across nodes and ensures
   that they are running and the choas workers in turn verify that the Minio
   server is running on the remote node at the specified port during initialization.
 - On the event of any chaos worker not reachable or if the Minio
   process is not running on the remote node the test fails.
 - Currrently RoundRobin Chaos Test is introduced wherein Minio server is
   stopped on each node and recovered after the specified recovery time one after the other.
   
# How to run. 

- Fetch the project 

  `$go get github.com/hackintoshrao/minio-chaos` 
- Build master program.

  `$ cd $GOPATH/src/github.com/hackintoshrao/minio-chaos/master && go build`
- Run Minio Distributed server using systemd script on remote nodes. 
  The workers use `systemd` to control the Minio process. [Click here](https://github.com/minio/minio/tree/master/dist/linux-systemd/distributed) for info on configuring systemd to run Minio Distributed.
- Build and run chaos-workers on each these remote nodes. Use `sudo` to run worker. 
  Need privilaged access to control Minio process using systemd.
  
  `$ go get github.com/hackintoshrao/minio-chaos`
  
  `$ cd $GOPATH/src/github.com/hackintoshrao/minio-chaos/worker/ && go build`
  
  `$ sudo ./worker`
- Run master. 
  `master -endpoints="<Node-1-IP>:9997,<NODE-2-IP>:9997,<NODE-3-IP>:9997...... -recover=30"`
   Currenly Chaos workers run at port 9997.
   
  
  

# Options

- `-recover : Removery time after the failure injection on remote Minio node`.

- `-endpoints: "," separted <IP>:<PORT> at which Remote Chaos Workers are Running"`.


# What next ?

- By Extending the `Chaos` interface different failures can be introduced. Next step is to make it more easier to run 
  , with more options and different failures like multiple node failures will introduced.
  
# Contributing. 

- Feel free to open and issue or send a PR.


