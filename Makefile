default: install

build:
	cd cmd/metricswatcher; go build

install:
	cd cmd/metricswatcher; go install

clean:
	rm cmd/metricswatcher/metricswatcher
