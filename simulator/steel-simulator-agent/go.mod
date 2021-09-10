module steel-simulator-agent

go 1.16

replace steel-simulator-common => ../steel-simulator-common

require (
	github.com/abu-lang/goabu v0.0.0
	steel-simulator-common v0.0.0
)

replace github.com/abu-lang/goabu => ../../goabu
