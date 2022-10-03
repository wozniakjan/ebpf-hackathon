// +build ignore

#include "common.h"

#include "bpf_tracing.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct event {
	u32 pid;
	u8 line[80];
};

struct {
	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} events SEC(".maps");

// Force emitting struct event into the ELF.
const struct event *unused __attribute__((unused));

SEC("uprobe/testbin_test")
int uprobe_testbin_test(struct pt_regs *ctx) {
    bpf_printk("Hello World  !\\n");
	return 0;
}
