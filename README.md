# scx-adapt

## what is it?

`scx-adapt` is basically a profiling tool for ***sched_ext*** schedulers. Main goal is to **choose and attach the appropiate scheduler** for the current workload type since there is not a **one-size-fits-all** solution for CPU scheduling.  

In YAML configuration files, ***sched_ext*** schedulers with their paths, priorities (scx-adapt iterates over schedulers and their criteria **in order of their priorities**) and their selection criteria are set. Schedulers get attached to kernel depending on the value of system performance metrics (e.g. load average). 

## example configuration file

```yaml
interval: 1000
schedulers:
  - path: "/etc/scx-adapt/schedulers/rr.o"
    priority: 1
    criterias:
      - value_name: load_avg_1
        more_than: 10
  - path: "/etc/scx-adapt/schedulers/cfs.o"
    priority: 2
    criterias:
      - value_name: io_psi_some_10
        more_than: 50
  - path: "/etc/scx-adapt/schedulers/lottery.o"
    priority: 3
    criterias:
      - value_name: load_avg_5
        less_than: 3 

```

- **interval**: System metrics reading period of scx-adapt (*in milliseconds*)
- **path**: Filepath of sched_ext scheduler (BPF bytecode). **Relative filepaths** are also supported but it is recommended to write **absolute paths**. (To prevent filepath based errors when added to `/etc/scx-adapt/schedulers`)
- **priority**: Priority of scheduler (1-139). scx-adapt starts checking the schedulers' criteria in order of their priorities (smaller value, higher priority). Attaches the first matching scheduler to the kernel.
- **criterias**: Criterias depending on the current value of system performance metrics. As of v0.0.1, it supports `more_than` and `less_than`
- **value_name**: System metric name
    - Currently supported metrics (*As of v0.0.1*):
        - (cpu|io|mem)\_psi_(some|full)_(10|60|300)
		- load_avg_(1|5|15)
		- procs_running
		- procs_blocked
		- procs_disk_io

## using scx-adapt

- `scx-adapt [command]`

- `add-profile <profile_path>` | Add scx-adapt profile configuration to profiles folder ('/etc/scx-adapt/profiles' by default)
- `check-profile` | Check if profile file in YAML format passed from STDIN is valid
- `install-service` | Add Systemd service file 'scx-adapt@.service' to '/etc/systemd/system'
- `list-profiles` | List profile configurations in profiles folder ('/etc/scx-adapt/profiles' by default)
- `log-csv <csv_file_path> [interval]` | Print system variables to file in csv format
- `log-sched` | Print sched_ext event tracing to stdout
- `remove-profile <profile_filename>` | Remove profile configuration from profiles folder ('/etc/scx-adapt/profiles' by default)
- `remove-service`  Remove Systemd service file 'scx-adapt@.service' in '/etc/systemd/system'
- `start-profile <profile_path>` | Run scx-adapt with the profile configuration
    - To start added profiles, pass profile filename (no need to specify directory, since scx-adapt scans '*/etc/scx-adapt/schedulers*' first)
- `status` | Print currently running sched_ext scheduler.
- `stop` | Stop currently running sched_ext scheduler

### Systemd service

To install the service file and enable/start the service:
- `$ scx-adapt install-service`
- `$ systemctl enable scx-adapt@<profile_path>`
- `$ systemctl start scx-adapt@<profile_path>`

To disable/stop and delete service file:
 
- `$ systemctl disable scx-adapt@<profile_path>`
- `$ systemctl stop scx-adapt@<profile_path>`
- `$ scx-adapt remove-service`


## installation

### go install

`go install github.com/dogukanmeral/scx-adapt@latest`

> Pre-compiled executables and distribution packages are not available as of v0.0.1 but I plan to add automated releases as soon as possible. 

### dependencies

- **bpftool**: Used while attaching the scheduler to the kernel.
- Kernel compiled with **BPF support** and BPF filesystem is mounted
- Kernel compiled with **sched_ext support** (sched_ext has been merged at Linux 6.12)

## further development

I (Doğukan Meral) have been the sole developer for v0.0.1 while my friend @onurkaragur is currently working on performance analysis of schedulers and ways to optimize scx-adapt configurations using machine learning techniques on the ***scx-adapt-experiments*** repository.

Your feedbacks, suggestions, criticisms and most importantly your contributions are highly appriciated. Feel free to contact me at my e-mail address `dogukan.meral@yahoo.com`   

## helpful resources for bpf and sched_ext / inspirations for this project

- Johannes Bechberger's [minimal scheduler repository](https://github.com/parttimenerd/minimal-scheduler) and [blog article](https://mostlynerdless.de/blog/2024/10/25/a-minimal-scheduler-with-ebpf-sched_ext-and-c/)
- Andrea Righi's [neural network scheduler video](https://youtu.be/ywW83YK9EKQ)
- [scx repository](https://github.com/sched-ext/scx) which contains many Sched_ext schedulers and tools
- [Perfetto](https://ui.perfetto.dev/): Browser based and locally running scheduler trace visualisation and analysis tool
- Changwoo Min's ['sched_ext: a BPF-extensible scheduler class'](https://blogs.igalia.com/changwoo/sched-ext-a-bpf-extensible-scheduler-class-part-1/) and ['sched_ext: scheduler architecture and interfaces'](blogs.igalia.com/changwoo/sched-ext-scheduler-architecture-and-interfaces-part-2/) blog articles

## license

scx-adapt is licensed under GPLv2. See LICENSE.