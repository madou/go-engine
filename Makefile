COVER_DIR = cover

coverage:
	mkdir -p $(COVER_DIR)
	go test github.com/walesey/go-engine/controller -coverprofile=$(COVER_DIR)/controller.cover.out && \
	go test github.com/walesey/go-engine/effects -coverprofile=$(COVER_DIR)/effects.cover.out && \
	go test github.com/walesey/go-engine/physics -coverprofile=$(COVER_DIR)/physics.cover.out && \
	go test github.com/walesey/go-engine/physics/gjk -coverprofile=$(COVER_DIR)/gjk.cover.out && \
	go test github.com/walesey/go-engine/util -coverprofile=$(COVER_DIR)/util.cover.out && \
	go test github.com/walesey/go-engine/vectormath -coverprofile=$(COVER_DIR)/vectormath.cover.out && \
		rm -f $(COVER_DIR)/coverage.out && \
		echo 'mode: set' > $(COVER_DIR)/coverage.out && \
		cat $(COVER_DIR)/*.cover.out | sed '/mode: set/d' >> $(COVER_DIR)/coverage.out && \
		go tool cover -html=$(COVER_DIR)/coverage.out -o=$(COVER_DIR)/coverage.html && \
		rm $(COVER_DIR)/*.cover.out