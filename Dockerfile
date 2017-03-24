FROM golang:1.7.5-wheezy

COPY . /home

WORKDIR /home

# supervisord.conf contains configuration for controlling Minio and worker processes.
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN \
       apt-get  --assume-yes update && \
       # install pip, required for installing supervisord.
       apt-get -y --assume-yes install python-pip && \
       # install supervisord.
       pip install supervisor

RUN \
       # Fetch and build  Minio server.
       go get -u github.com/minio/minio && \
       # Build worker from the local changes.
       go build /home/worker/chaos-worker.go && \
       # Fetch and build mc. 
       go get -u github.com/minio/mc && \
       # configure mc to point to the Minio server running in the container.
       $GOPATH/bin/mc config host add myminio http://127.0.0.1:9000 minio123 minio12345 && \
       # Create log dir.
       mkdir /var/log/minio/ && \
       # supervisord needs the binaries to be in /usr/local/bin/
       mv $GOPATH/bin/minio /usr/local/bin/ && \
       mv /home/chaos-worker /usr/local/bin/worker && \
       chmod +x /home/start.sh && \
       sync

# supervisord manages Minio and worker processses.
CMD /home/start.sh ; /usr/local/bin/supervisord -c /home/supervisord.conf 
