# docker image registry
IMG := wozniakjan

# compile eBPF program using cillium/ebpf library 
bpf_bpfel_x86.o: main.go
	CPATH=../headers go generate ./...

# testbin is a sample application used for inspection by tracer
testbin: test_bin.go
	# turn of inlining and optimizations
	go build -gcflags '-N -l' -o testbin ./test_bin.go

# tracer is an agent application that will instrument testbin by uprobe on a function
tracer: main.go bpf_bpfel_x86.o
	go build -o tracer ./main.go ./bpf_bpfel_x86.go

# build docker image
image-%:
	docker build . -t $(IMG)/$* -f deploy/Dockerfile.$*

.PHONY: images
images: testbin tracer image-testbin image-tracer

# push docker image
push-%:
	docker push $(IMG)/$*

.PHONY: push-images
push-images: push-testbin push-tracer

# run both testbin and tracer locally with default config
.PHONY: run
run: tracer testbin
	./testbin > /dev/null &
	sudo ./tracer	

.PHONY: clean
clean:
	rm -rf bpf_bpfel_x86.go bpf_bpfel_x86.o testbin tracer
