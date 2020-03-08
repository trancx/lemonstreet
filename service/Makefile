PROJECT_NAME	:=  lemonstreet
DOCKER_VOLUME_DIR := /mnt
DB_VOLUME	:= $(DOCKER_VOLUME_DIR)/$(PROJECT_NAME)/db

DISC_VOLUME	:= $(DOCKER_VOLUME_DIR)/discovery
DISCOVERY	:= 

NGX_VOLUME	:= $(DOCKER_VOLUME_DIR)/nginx


service_dirs	:=  account/ article/ comment/ explore/ login/ search/ verify/
services	:=  $(patsubst %/, %, $(service_dirs))
targets		:=  $(patsubst %/, %/cmd/cmd, $(service_dirs))	
dtargets	:=  $(patsubst %/, $(DOCKER_VOLUME_DIR)/$(PROJECT_NAME)/%, $(service_dirs))

PHONY := all
Q := @
CGO_ENABLED := 0

export CGO_ENABLED
export Q
export PROJECT_NAME
export DOCKER_VOLUME_DIR
# go code -> executable -> dockerfile -> image

all: $(services) $(DB_VOLUME) $(DISC_VOLUME) $(NGX_VOLUME)
	

$(services): FORCE
	$(Q)$(MAKE) -C $@ all
	$(Q)chmod 777 $(targets)
	
# build DB
$(DB_VOLUME):
	$(Q)mkdir -p $(DB_VOLUME)/initdb.d $(DB_VOLUME)/logs $(DB_VOLUME)/data
	$(Q)docker volume create --driver local --opt type=none --opt device=$(DB_VOLUME)/initdb.d --opt o=bind --name=db-init-volume
	$(Q)docker volume create --driver local --opt type=none --opt device=$(DB_VOLUME)/data --opt o=bind --name=db-data-volume
	$(Q)docker volume create --driver local --opt type=none --opt device=$(DB_VOLUME)/logs --opt o=bind --name=db-log-volume

$(DISC_VOLUME):
	$(Q)mkdir -p $(DISC_VOLUME)
	$(Q)docker volume create --driver local --opt type=none --opt device=$(DISC_VOLUME) --opt o=bind --name=disc-data-volume
	$(Q)cp $(DISCOVERY)/cmd $(DISCOVERY)/*.toml -t $(DISC_VOLUME)

$(NGX_VOLUME):
	$(Q)mkdir -p $(NGX_VOLUME)/conf $(NGX_VOLUME)/logs $(NGX_VOLUME)/res
	$(Q)docker volume create --driver local --opt type=none --opt device=$(NGX_VOLUME)/res --opt o=bind --name=ngx-res-volume
	$(Q)docker volume create --driver local --opt type=none --opt device=$(NGX_VOLUME)/logs --opt o=bind --name=ngx-log-volume
	$(Q)docker volume create --driver local --opt type=none --opt device=$(NGX_VOLUME)/conf --opt o=bind --name=ngx-conf-volume

PHONY += clean
clean:
	$(Q)rm -rf $(targets)

PHONY += d_clean
d_clean:
	$(Q)rm -rf $(dtargets)

PHONY += FORCE
FORCE:

.PHONY: $(PHONY)
