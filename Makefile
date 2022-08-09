BUILD=build
PY = python3.6 -O -m compileall -b -q -f
PYSRC = demo/

all: clean build pydemo

build: go.sum
	@echo "Building ..."
	@go build -mod=readonly -o $(BUILD)/multinfer

go.sum: go.mod
	@echo "Ensure dependencies have not been modified"
	@GO111MODULE=on go mod verify


pydemo:
	@echo "Compiling demo ..."
	@mkdir -p $(BUILD)
	@cp -r $(PYSRC) $(BUILD)/
	@-$(PY) $(BUILD)/demo
	@find $(BUILD)/demo -name '*.py' -delete
	@find $(BUILD)/demo -name "__pycache__" |xargs rm -rf
	@rm $(BUILD)/demo/config/settings.pyc
	@cp $(PYSRC)/config/settings.py $(BUILD)/demo/config/

clean:
	@echo "Clean old built"
	@rm -rf $(BUILD)
	@go clean
	@find . -name "__pycache__" | xargs rm -rf
	@find . -name '*.pyc' -delete

