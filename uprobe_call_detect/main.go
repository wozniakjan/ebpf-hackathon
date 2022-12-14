// This is an implementation of sample agent application that injects ebpf program as
// a hook on a certain binary and function

//go:build linux
// +build linux

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/rlimit"
)

// generate bpf_bpfel_x86.go and compile bpf_bpfel_x86.o
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" -target native -type event bpf uprobe.c -- -I../headers

const (
	symbol = "main.easyToFindFunctionName"
)

func tracedBinPath(bin *string) []string {
	if bin != nil && *bin != "" {
		return strings.Split(*bin, ":")
	}
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	return []string{path.Join(exPath, "testbin")}
}

func instrumentBin(binPath string, objs bpfObjects) func() error {
	if binPath == "" {
		return nil
	}
	log.Println("instrumenting", binPath)
	ex, err := link.OpenExecutable(binPath)
	if err != nil {
		log.Fatalf("opening executable: %s", err)
	}

	// Open a Uprobe at the entry point of the symbol and attach the pre-compiled eBPF program to it.
	up, err := ex.Uprobe(symbol, objs.UprobeTestbinTest, nil)
	if err != nil {
		log.Fatalf("creating uprobe: %s", err)
	}
	return up.Close
}

func main() {
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	bin := flag.String("bin", "", "colon separated list of paths to the instrumented binaries")
	flag.Parse()

	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	// Load pre-compiled programs and maps into the kernel.
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %s", err)
	}
	defer objs.Close()

	// Open an ELF binary and read its symbols.
	binPaths := tracedBinPath(bin)
	for _, binPath := range binPaths {
		cl := instrumentBin(binPath, objs)
		defer cl()
	}

	// Open a perf event reader from userspace on the PERF_EVENT_ARRAY map described in the eBPF C program.
	rd, err := perf.NewReader(objs.Events, os.Getpagesize())
	if err != nil {
		log.Fatalf("creating perf event reader: %s", err)
	}
	defer rd.Close()

	go func() {
		// Wait for a signal and close the perf reader, which will interrupt rd.Read() and make the program exit.
		<-stopper
		log.Println("Received signal, exiting program..")

		if err := rd.Close(); err != nil {
			log.Fatalf("closing perf event reader: %s", err)
		}
	}()

	log.Printf("Listening for events..")

	// bpfEvent is generated by bpf2go.
	var event bpfEvent
	for {
		record, err := rd.Read()
		if err != nil {
			if errors.Is(err, perf.ErrClosed) {
				return
			}
			log.Printf("reading from perf event reader: %s", err)
			continue
		}

		if record.LostSamples != 0 {
			log.Printf("perf event ring buffer full, dropped %d samples", record.LostSamples)
			continue
		}

		// Parse the perf event entry into a bpfEvent structure.
		if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("parsing perf event: %s", err)
			continue
		}

		// TODO: add bin path support for multiple paths
		log.Printf("%v %s:%s argument: %v %v", event.Pid, "binPath", symbol, event.Arg, event.Ret)
	}
}
