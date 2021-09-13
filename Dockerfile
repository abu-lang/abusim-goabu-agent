# Create the building image for compiling
FROM golang:1.16-alpine as build

RUN mkdir /agent
WORKDIR /agent

COPY ./abusim-core-dev/schema ./abusim-core-dev/schema
COPY ./goabu ./goabu

WORKDIR /agent/abusim-goabu-agent

COPY ./abusim-goabu-agent-dev/go.mod .
COPY ./abusim-goabu-agent-dev/go.sum .
RUN go mod download -x

COPY ./abusim-goabu-agent-dev .

RUN CGO_ENABLED=0 go build

# Create the final image with the executable
FROM scratch as exec

COPY --from=build /agent/abusim-goabu-agent/abusim-goabu-agent /agent

ENTRYPOINT [ "/agent" ]
