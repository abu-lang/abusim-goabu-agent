module abusim-goabu-agent

go 1.16

require (
	github.com/abu-lang/abusim-core/schema v0.0.0
	github.com/abu-lang/goabu v0.0.0
)

replace github.com/abu-lang/abusim-core/schema => ../abusim-core-dev/schema
replace github.com/abu-lang/goabu => ../goabu
