#!/bin/sh
# Generate the vmlinux.h file to start using in the IDE

bpftool btf dump file /sys/kernel/btf/vmlinux format c > sample-bpf/vmlinux.h
