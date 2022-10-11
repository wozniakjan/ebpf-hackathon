bpf_bpfel_x86.o: main.go
	CPATH=../headers go generate ./...

testbin: test_bin.go
	# turn of inlining and optimizations
	go build -gcflags '-N -l' -o testbin ./test_bin.go

tracer: main.go bpf_bpfel_x86.o
	go build -o tracer ./main.go ./bpf_bpfel_x86.go

.PHONY: run
run: tracer testbin
	sudo ./tracer	

.PHONY: clean
clean:
	rm -rf bpf_bpfel_x86.go bpf_bpfel_x86.o testbin tracer
