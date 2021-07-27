# Create the building image for compiling
FROM golang:1.16-alpine as build

RUN mkdir /agent
WORKDIR /agent

COPY ./simulator/steel-simulator-common ./simulator/steel-simulator-common
COPY ./src ./src

WORKDIR /agent/simulator/steel-simulator-agent

COPY ./simulator/steel-simulator-agent/go.mod .
COPY ./simulator/steel-simulator-agent/go.sum .
RUN go mod download -x

COPY ./simulator/steel-simulator-agent .

RUN CGO_ENABLED=0 go build

# Create the final image with the executable
FROM scratch as exec

COPY --from=build /agent/simulator/steel-simulator-agent/steel-simulator-agent /agent

ENTRYPOINT [ "/agent" ]
