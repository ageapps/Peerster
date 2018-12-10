t = 3 # timer
n = 4 # node number
f = test.png
h = 13fa82c9e76e18e1e8587231be1aa955f3469a20a1b085a28326339f36108ddd
d = nodeC

build:
	go build .
	cd ./client && go build . && cd ..
	cd ./server && go build . && cd ..

clean:
	go clean
	go clean ./client/
	go clean ./server/

run:
	go run --race . -UIPort=1000$(n) -gossipAddr=127.0.0.1:500$(n) -name=node$(n) -rtimer=$(t)

run1:
	go run --race . -UIPort=10000 -gossipAddr=127.0.0.1:5000 -name=nodeA -rtimer=3

run2:
	go run --race . -UIPort=10001 -gossipAddr=127.0.0.1:5001 -name=nodeB -peers=127.0.0.1:5000 -rtimer=3

run3:
	go run --race . -UIPort=10002 -gossipAddr=127.0.0.1:5002 -name=nodeC -peers=127.0.0.1:5001 -rtimer=3

send:
	go run --race ./client -UIPort=10001 -msg=Hello -Dest=$(d) -file=$(f) -request=$(h)

serve:
	cd ./server && go run --race .

private:
	go run --race ./client -UIPort=10002 -msg=Hello -Dest=$(d)
	
front:	
	location=~/git/Peerster-App; \
	current=$(shell pwd) && cd $$location && npm run build && cd $$current; \
	bash -c "rm -r web/*"; \
	cp -R $$location/dist/* ./web 

test1:
	sh test/test_1_ring.sh

test2:
	sh test/test_2_ring.sh

cchunks:
	rm ./._Chunks/*
	rm ./._Metafiles/*
