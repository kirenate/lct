FROM golang:1.24-bullseye AS build

WORKDIR /src

# Copy only dependency files first
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

RUN cd backend &&  GO111MODULE=on CGO_ENABLED=1 go build -o=bootstrap main.go

## Now copy it into our base image.
FROM debian:trixie

RUN rm -rf /var/lib/apt/lists/* \
    && apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl
RUN update-ca-certificates
RUN rm -rf /var/lib/apt/lists/*

COPY --from=build /src/backend/bootstrap /

CMD ["/bootstrap"]