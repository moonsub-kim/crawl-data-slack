FROM golang:1.17.1 AS build

# Change working dir -> for locating configrations
RUN mkdir -p /app
WORKDIR /app

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# Copy the local package files to the container's workspace.
COPY . .

# Build the outyet command inside the container.
RUN go build -o /go/bin/crawl-data-slack ./cmd

FROM golang:1.17.1 AS prod

COPY --from=build /go/bin/crawl-data-slack /go/bin/crawl-data-slack

COPY --from=build /app/wait-for-it.sh /bin/wait-for-it.sh
