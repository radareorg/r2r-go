BINFOLDER := bin
R2RMAIN := r2r
R2RBUILDER := r2r-build

all: setup main
	@echo "Built."

clean:
	rm -rf $(BINFOLDER)

setup:
	@echo "[MKDIR]" $(BINFOLDER)
	@mkdir -p $(BINFOLDER)

main:
	@echo "[GO]" $(R2RMAIN)
	@cd $(R2RMAIN); go get r2pipe; go build
	@echo "[MV]" $(R2RMAIN)"/r2r"
	@mv $(R2RMAIN)/r2r $(BINFOLDER)/r2r

builder: setup
	@echo "[GO]" $(R2RBUILDER)
	@cd $(R2RBUILDER); go build
	@echo "[MV]" $(R2RBUILDER)"/r2r"
	@mv $(R2RBUILDER)/r2r-build $(BINFOLDER)/r2r-build

