BINFOLDER := $(CURDIR)/bin
GOPATH := $(CURDIR)/gopath
R2RMAIN := r2r
R2RBUILDER := r2r-build

all: setup main
	@echo "Built."

clean:
	rm -rf $(BINFOLDER) $(GOPATH)

setup: $(BINFOLDER) $(GOPATH)

$(BINFOLDER):
	@echo "[MKDIR]" $(BINFOLDER)
	@mkdir -p $(BINFOLDER)

$(GOPATH):
	@echo "[MKDIR]" $(GOPATH)
	@mkdir -p $(GOPATH)

main:
	@echo "[GO]" $(R2RMAIN)
	@cd $(R2RMAIN); export GOPATH=$(GOPATH) ; go get github.com/radare/r2pipe-go
	@cd $(R2RMAIN); export GOPATH=$(GOPATH) ; go build
	@echo "[MV]" $(R2RMAIN)"/r2r"
	@mv $(R2RMAIN)/r2r $(BINFOLDER)/r2r

builder: setup
	@echo "[GO]" $(R2RBUILDER)
	@cd $(R2RBUILDER); export GOPATH=$(GOPATH) ; go build
	@echo "[MV]" $(R2RBUILDER)"/r2r"
	@mv $(R2RBUILDER)/r2r-build $(BINFOLDER)/r2r-build

