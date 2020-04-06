target	:=	cmd/cmd
PHONY	:= 	all
service :=  	search
ifeq ($(target), $(wildcard $(target)))
update := $(shell find . -name '*.go' -newer cmd/cmd | sed 's/\.\///g' | grep go)
endif

ifdef PROJECT_NAME
ifdef DOCKER_VOLUME_DIR
volume_dir := $(DOCKER_VOLUME_DIR)/$(PROJECT_NAME)/$(service)
dtarget := $(DOCKER_VOLUME_DIR)/$(PROJECT_NAME)/$(service)/cmd

all: $(dtarget)

$(dtarget): $(target) $(volume_dir)
	$(Q)cp -fr $(target) configs/ -t $(volume_dir)

ifeq ($(update),) 
$(target):
	$(Q)cd cmd && go build && echo "rebuild the go service"
else
$(target): FORCE
	$(Q)cd cmd && go build && echo "rebuild the go service"
endif #update
$(volume_dir):
	$(Q)mkdir -p $(volume_dir)
	$(Q)docker volume create --driver local --opt type=none --opt device=$(volume_dir) --opt o=bind  --name=$(service)-data-volume
	$(Q)cp wait-for-it.sh -t $(volume_dir)
endif #DOCKER_VOLUME_DIR
else
all:
	$(Q)echo "please make in project directory"
	$(Q)exit 1
endif #PROJECT_NAME

PHONY += clean
clean:
	$(Q)rm $(target)

PHONEY += FORCE
FORCE:

.PHONEY: $(PHONEY)
