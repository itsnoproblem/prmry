build:
	go build -o server ./cmd/htmx-server/

run:
	make build
	./server

watch:
	ulimit -n 9999 #increase the file watch limit, might required on MacOS
	reflex -s -r '\.(go|gohtml)$$' make run
	#reflex -s -g './pkg/templates/*.gohtml' make run
