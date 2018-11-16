build:
	go build .
	cd ./client && go build . && cd ..
	cd ./server && go build . && cd ..

clean:
	go clean
	go clean ./client/
	go clean ./server/

test1:
	go run --race . -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -rtimer=3

test2:
	go run --race . -UIPort=10001 -gossipAddr=127.0.0.1:5001 -name=nodeB -peers=127.0.0.1:5000 -rtimer=3

test3:
	go run --race . -UIPort=10002 -gossipAddr=127.0.0.1:5002 -name=nodeC -peers=127.0.0.1:5001 -rtimer=3

send:
	go run --race ./client -UIPort=10000 -msg=Hello

serve:
	go run --race ./server

private:
	go run --race ./client -UIPort=10002 -msg=Hello -dest=nodeA
	
front:	
	bash -c "rm -r app/*"; \
	location=~/git/Peerster-App; \
	current=$(shell pwd) && cd $$location && npm run build && cd $$current; \
	cp -R $$location/dist/* ./app 

	

	
  
