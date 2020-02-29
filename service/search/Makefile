target	:=	cmd/cmd
PHONY	:= 	all
update	:=	$(shell find . -name '*.go' -newer cmd/cmd | sed 's/\.\///g' | grep go)

ifeq ($(update),)
#no need to rebuild
all:
	$(Q)echo "every thing is up to date"
else
all:  $(target)
endif

$(target): FORCE
	$(Q)cd cmd/ && go build && echo "rebuild the go service"

PHONEY += FORCE
FORCE:

.PHONEY: $(PHONEY)
