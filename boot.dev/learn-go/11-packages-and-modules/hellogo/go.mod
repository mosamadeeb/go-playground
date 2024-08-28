module github.com/mosamadeeb/hellogo

go 1.23.0

// This basically overrides the remote package with a local one
replace github.com/mosamadeeb/mystrings v0.0.0 => ../mystrings

require (
    github.com/mosamadeeb/mystrings v0.0.0
)
