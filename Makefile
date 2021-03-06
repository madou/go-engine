COVER_DIR = cover

testintall:
	go get -t github.com/walesey/go-engine/actor
	go get -t github.com/walesey/go-engine/assets
	go get -t github.com/walesey/go-engine/controller
	go get -t github.com/walesey/go-engine/effects
	go get -t github.com/walesey/go-engine/engine
	go get -t github.com/walesey/go-engine/networking
	go get -t github.com/walesey/go-engine/physics/physicsAPI
	go get -t github.com/walesey/go-engine/ui
	go get -t github.com/walesey/go-engine/util

test: testintall
	go test github.com/walesey/go-engine/actor
	go test github.com/walesey/go-engine/assets
	go test github.com/walesey/go-engine/controller
	go test github.com/walesey/go-engine/effects
	go test github.com/walesey/go-engine/engine
	go test github.com/walesey/go-engine/networking
	go test github.com/walesey/go-engine/physics/physicsAPI
	go test github.com/walesey/go-engine/ui
	go test github.com/walesey/go-engine/util

coverage: testintall
	mkdir -p $(COVER_DIR)
	go test github.com/walesey/go-engine/actor -coverprofile=$(COVER_DIR)/actor.cover.out && \
	go test github.com/walesey/go-engine/assets -coverprofile=$(COVER_DIR)/assets.cover.out && \
	go test github.com/walesey/go-engine/controller -coverprofile=$(COVER_DIR)/controller.cover.out && \
	go test github.com/walesey/go-engine/effects -coverprofile=$(COVER_DIR)/effects.cover.out && \
	go test github.com/walesey/go-engine/engine -coverprofile=$(COVER_DIR)/engine.cover.out && \
	go test github.com/walesey/go-engine/networking -coverprofile=$(COVER_DIR)/networking.cover.out && \
	go test github.com/walesey/go-engine/physics/physicsAPI -coverprofile=$(COVER_DIR)/physics.cover.out && \
	go test github.com/walesey/go-engine/ui -coverprofile=$(COVER_DIR)/ui.cover.out && \
	go test github.com/walesey/go-engine/util -coverprofile=$(COVER_DIR)/util.cover.out && \
		rm -f $(COVER_DIR)/coverage.out && \
		echo 'mode: set' > $(COVER_DIR)/coverage.out && \
		cat $(COVER_DIR)/*.cover.out | sed '/mode: set/d' >> $(COVER_DIR)/coverage.out && \
		go tool cover -html=$(COVER_DIR)/coverage.out -o=$(COVER_DIR)/coverage.html && \
		rm $(COVER_DIR)/*.cover.out
