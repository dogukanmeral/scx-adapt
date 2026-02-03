#!/bin/sh

# Usage: ./build.sh bpf_file.c (default is sched_ext.bpf.c)

# Set the default file
BPF_FILE=${1:-sample-bpf/sched_ext.bpf.c}
BASE_FILE=$(basename ${BPF_FILE})


# Create the vmlinux header with all the eBPF Linux functions
# if it doesn'r exist
if [ ! -f sample-bpf/vmlinux.h ]; then
    echo "Creating vmlinux.h"
    bpftool btf dump file /sys/kernel/btf/vmlinux format c > sample-bpf/vmlinux.h
fi

# Compile the scheduler
clang -target bpf -g -O2 -c $BPF_FILE -o bytecode/${BASE_FILE}.o -I.