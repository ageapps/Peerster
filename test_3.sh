#!/usr/bin/env bash

go build
cd client
go build
cd ..

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
DEBUG="true"

outputFiles=()
message_c1_1=Weather_is_clear
message_c2_1=Winter_is_coming
message_c1_2=No_clouds_really
message_c2_2=Let\'s_go_skiing
message_c3=Is_anybody_here?


UIPort=12345
gossipPort=5000
name='A'

# General peerster (gossiper) command
#./Peerster -UIPort=12345 -gossipAddr=127.0.0.1:5001 -name=A -peers=127.0.0.1:5002 > A.out &

for i in `seq 1 10`;
do
	outFileName="$name.out"
	peerPort=$((($gossipPort+1)%10+5000))
	peer="127.0.0.1:$peerPort"
	gossipAddr="127.0.0.1:$gossipPort"
	go run -race main.go -UIPort=$UIPort -gossipAddr=$gossipAddr -name=$name -peers=$peer > $outFileName &
	outputFiles+=("$outFileName")
	if [[ "$DEBUG" == "true" ]] ; then
		echo "$name running at UIPort $UIPort and gossipPort $gossipPort"
	fi
	UIPort=$(($UIPort+1))
	gossipPort=$(($gossipPort+1))
	name=$(echo "$name" | tr "A-Y" "B-Z")
done
sleep 5

# E
./client/client -UIPort=12349 -msg=$message_c1_1
# B
./client/client -UIPort=12346 -msg=$message_c2_1
sleep 2
# E
./client/client -UIPort=12349 -msg=$message_c1_2
# B
./client/client -UIPort=12346 -msg=$message_c2_2
sleep 1
# G
./client/client -UIPort=12351 -msg=$message_c3
sleep 5
pkill -f go

