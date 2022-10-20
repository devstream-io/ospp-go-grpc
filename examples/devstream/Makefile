.PHONY: run
run:
	go build -o core ./core
	go build -o plugin ./plugin
	chmod +x ./core/core
	chmod +x ./plugin/plugin
	cd core && ./core

.PHONY: delete
delete:
	rm -f ./core/core && rm -f ./plugin/plugin
