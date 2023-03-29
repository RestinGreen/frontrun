# To run the project without docker
clear && go build -o main ./cmd/main/ && ./main

# To run the project with docker

1. Build + Run
    docker rm frontrun && docker rmi frontrun && docker build . -t frontrun

2. Run ( -v param needs to be the location of the geth IPC on the Host machine ) 

    - on mainnet (-v: pass the location of the ipc)
   docker run -v /root/eth/execution/chaindata/:/app/ipc/ --env TESTNET=false -d --name frontrun frontrun

    - on testnet (-v pass the location of the ipc)
    docker run -v /root/goerli/execution/chaindata/:/app/ipc/ --env TESTNET=true -d --name frontrun  frontrun

# Docker commands cheat sheet

- show logs of running container
    docker logs <name> --tail 20 -f 
- show list of all running and stopped containers
    docker ps -a
- list all docker images
    docker images -a
