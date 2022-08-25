module github.com/damejeras/ferry/example

go 1.18

require github.com/damejeras/ferry v0.0.0-20220812150200-855f29bea2a6

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
)

require (
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
)

replace github.com/damejeras/ferry => ../.
