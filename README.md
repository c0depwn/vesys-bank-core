# VESYS 2023 bank-server

This project is a backend for the FHNW course VESYS.
It provides server-side functionality for an existing GUI.

## Running in docker

```shell
# start the container, the API should become available at localhost:8080, for more details check the compose file
docker compose -f deployments/local/docker-compose.yaml up -d 

# stopping
docker compose -f deployments/local/docker-compose.yaml down
```

