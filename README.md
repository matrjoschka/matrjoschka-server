# Matrjoschka Server

## Local Development
1. Install and Setup go
1. Install Glide https://github.com/Masterminds/glide
1. Run Glide Install `glide install`
1. Start the Server `go run .`
1. The API is running on Port `9287`

## Publish the docker image
1. Install and Setup Docker
1. Login at DockerHub `docker login`
1. Build the Image `docker build .`
1. Tag the Image `docker tag <image hash from previous step> <your docker hub name>/matrjoschka-server:latest`
1. Push the Image `docker push`