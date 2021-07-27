module steel-simulator-agent

go 1.16

replace (
	steel-lang => ../../src
	steel-simulator-common => ../steel-simulator-common
)

require (
	steel-lang v0.0.0
	steel-simulator-common v0.0.0
)
