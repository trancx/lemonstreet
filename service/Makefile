service_dirs	:=  account/ article/ comment/ explore/ login/ search/ verify/
services	:=  $(patsubst %/, %, $(service_dirs))
targets		:=  $(patsubst %/, %/cmd/cmd, $(service_dirs))	
PHONY := all
Q := @

export Q
# go code -> executable -> dockerfile -> image

all: $(services)
	$(Q)echo ${services}


$(services): FORCE
	$(Q)echo "" $@ 
	$(Q)$(MAKE) -C $@ all

PHONY += FORCE
FORCE:

.PHONY: $(PHONY)
