// +build ignore

#include "common.h"

#include "bpf_tracing.h"

char __license[] SEC("license") = "Dual MIT/GPL";

struct event {
	u32 pid;
	u32 arg;
    long ret;
};

struct {
	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
} events SEC(".maps");

// Force emitting struct event into the ELF.
const struct event *unused __attribute__((unused));

SEC("uprobe/testbin_test")
int uprobe_testbin_test(struct pt_regs *ctx) {
    bpf_printk("%d %d", PT_REGS_PARM1(ctx));

    struct event event;

    // go gc passes all args through stack, so can't use PT_REGS_PARM1(ctx) to get values for args
    // and instead need to look back with Stack Pointer (plan9 convention)
    // cgo follows the standard AMD64 ABI passing through registers
    // event.arg = (u32)PT_REGS_PARM1(ctx);
    void* stackAddr = (void*)PT_REGS_SP(ctx);
    event.ret = bpf_probe_read(&event.arg, sizeof(event.arg), stackAddr+8);
    event.pid = bpf_get_current_pid_tgid();

    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &event, sizeof(event));

	return 0;
}
