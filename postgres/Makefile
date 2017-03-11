REPO = "ahouston"
IMAGE = "flight-pgsql"
VERSION = "0.1.0"

PORTS = "5432:5432"

build:
	docker build -t $(REPO)/$(IMAGE):$(VERSION) . 

shell:
	docker run --rm -it $(REPO)/$(IMAGE):$(VERSION) /bin/bash

run: 
	docker run --rm -d -p $(PORTS) $(REPO)/$(IMAGE):$(VERSION)
