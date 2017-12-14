export UDIR= .
export GOC = x86_64-xen-ethos-6g
export GOL = x86_64-xen-ethos-6l
export ETN2GO = etn2go
export ET2G   = et2g
export EG2GO  = eg2go

export GOARCH = amd64
export TARGET_ARCH = x86_64
export GOETHOSINCLUDE=/usr/lib64/go/pkg/ethos_$(GOARCH)
export GOLINUXINCLUDE=/usr/lib64/go/pkg/linux_$(GOARCH)

export ETHOSROOT=server/rootfs
export MINIMALTDROOT=server/minimaltdfs

.PHONY: all install

program = test
new_type = commandTypes
New_type = CommandList
new_type_index = index

all: $(program)

install: $(program)
	sudo rm -rf server/
	(ethosParams server && cd server/ && ethosBuilder && minimaltdBuilder)
	ethosTypeInstall $(new_type)
	ethosDirectoryInstall user/nobody/ $(ETHOSROOT)/types/spec/$(new_type)/$(New_type) all
	cp $(program) $(ETHOSROOT)/programs
	ethosStringEncode /programs/$(program) > $(ETHOSROOT)/etc/init/services/$(program)

$(new_type).go: commandTypes.t
	$(ETN2GO) . $(new_type) main $^

$(program): $(program).go $(new_type).go
	ethosGo $^ 

run:
	make
	sudo make install
	cd server/ ; sudo ethosRun -t
	cd server/rootfs/log/test ; ethosLog .

clean:
	rm -rf $(new_type)/ $(new_type_index)/
	rm -f $(new_type).go
	rm -f $(program)
	rm -f $(program).goo.ethos
