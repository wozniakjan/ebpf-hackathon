# Kubermatic Hackathon 2022

## Learn some eBPF for greater good

There is a trend in writing eBPF programs for cloud native infrastructure to achieve things that were previously considered impossible. This topic is aimed at eBPF novices interested in getting their feet wet by writing some useful eBPF program from scratch.

**Stretch goal:** prepare workshop that can be submitted to conferences or online to enable enhanced learning experience to developers who are new to kernel

## Repository content

* [`./headers`](./headers) - basic eBPF C headers for efficiently writing eBPF programs
* [`./uprobe_call_detect`](./uprobe_call_detect) - example with uprobe and function argument inspection
* [`./presentation.pdf`](./presentation.pdf) - presentation for hackathon
  * [@sachintiptur](https://github.com/sachintiptur) with ebpf introduction
  * [@wozniakjan](https://github.com/wozniakjan) with how to write uprobe example
  * [@hdurand0710](https://github.com/hdurand0710) and [@stroebitzer](https://github.com/stroebitzer) on their exploration of https://github.com/parca-dev/parca and https://github.com/sustainable-computing-io/kepler
