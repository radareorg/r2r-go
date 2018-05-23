## 
## Copyright (c) 2018, Giovanni Dante Grazioli <deroad@libero.it>
## All rights reserved.
## 
## Redistribution and use in source and binary forms, with or without
## modification, are permitted provided that the following conditions are met:
## 
## * Redistributions of source code must retain the above copyright notice, this
##   list of conditions and the following disclaimer.
## * Redistributions in binary form must reproduce the above copyright notice,
##   this list of conditions and the following disclaimer in the documentation
##   and/or other materials provided with the distribution.
## 
## THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
## AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
## IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
## ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
## LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
## CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
## SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
## INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
## CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
## ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
## POSSIBILITY OF SUCH DAMAGE.
## 

GOPATH     := $(CURDIR)/gopath
BINFOLDER  := $(CURDIR)/bin
GODIFF     := github.com/pmezard/go-difflib/difflib
GODIFFPATH := $(GOPATH)/src/$(GODIFF)
R2RMAIN    := r2r
R2RBUILDER := r2r-build
GO         := export GOPATH=$(GOPATH) ; go 

all: setup main
	@echo "Built."

clean:
	rm -rf $(BINFOLDER) $(GOPATH)

setup: $(BINFOLDER) $(GOPATH) $(GODIFFPATH)

$(BINFOLDER):
	@echo "[MKDIR]" $(BINFOLDER)
	@mkdir -p $(BINFOLDER)

$(GODIFFPATH):
	@echo "[DEPS] diff"
	@$(GO) get -u -v $(GODIFF)

$(GOPATH):
	@echo "[MKDIR]" $(GOPATH)
	@mkdir -p $(GOPATH)

main:
	@echo "[GO]" $(R2RMAIN)
	@cd $(R2RMAIN); $(GO) build
	@echo "[MV] r2r"
	@mv $(R2RMAIN)/r2r $(BINFOLDER)/r2r

builder: setup
	@echo "[GO]" $(R2RBUILDER)
	@cd $(R2RBUILDER); $(GO) build
	@echo "[MV] r2r-build"
	@mv $(R2RBUILDER)/r2r-build $(BINFOLDER)/r2r-build

