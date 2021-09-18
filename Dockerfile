# Create the building image for compiling
FROM golang:1.16-alpine as build

RUN apk update
RUN apk add git

RUN mkdir /agent

WORKDIR /agent/abusim-goabu-agent

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 go build

# Create the final image with the executable
FROM scratch as exec

COPY --from=build /agent/abusim-goabu-agent/abusim-goabu-agent /agent

ENTRYPOINT [ "/agent" ]
