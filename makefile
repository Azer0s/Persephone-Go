all:
	@$(MAKE) clear
	@$(MAKE) persephone

clear:
	@echo
	@echo "\033[4m\033[1mClearing output folder\033[0m"
	@echo
	@rm -rf bin/ || true
	@mkdir bin

persephone:
	@echo
	@echo "\033[4m\033[1mBuilding runtime\033[0m"
	@echo
	@go build -o bin/persephone src/main.go

run/%.psph: %.psph
	@./bin/persephone -i $<

test:
	@eval "for f in examples/*.psph; do ./bin/persephone -i $f; done"