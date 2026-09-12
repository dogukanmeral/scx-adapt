# scx-adapt

`scx-adapt` is basically an automatic selection tool for ***sched_ext*** schedulers.

Main goal is to **choose and attach the appropriate scheduler** for the current workload type since there is not a **one-size-fits-all** solution for CPU scheduling.

## Features

- YAML-based profile configuration
- Automatic scheduler selection based on the current workload
- Priority-based ordering of schedulers
- Support for `external` (BPF bytecode) and `builtin` (executable) schedulers
- Runtime parameters for `builtin` schedulers
- Selection criteria using live metrics: load average, PSI (pressure), process counts, disk I/O
- Fallback to the system scheduler when no criteria match
- BPF loading via `cilium/ebpf` (no `bpftool` dependency)
- Object-grouped CLI: `profile`, `scheduler`, `service`
- Optional per-scheduler logging to file
- sched_ext event tracing
- Live status overview

---

## Profiles

Example:
```yaml
interval: 1000
schedulers:
  - path: "rr.o"
    loader: external
    structName: "sched_ops"
    priority: 1
    criterias:
      - value_name: load_avg_1
        less_than: 2
  - path: "/home/dogukan/schedulers/cfs.o"
    loader: external
    structName: "sched_ops"
    priority: 2
    criterias:
      - value_name: io_psi_some_10
        more_than: 50
  - path: "scx_central"
    loader: builtin
    log: true
    parameters:
      - "-s=50"
    priority: 3
    criterias:
      - value_name: load_avg_5
        more_than: 4

```

| Key | Required | Applies to | Description |
| --- | --- | --- | --- |
| `interval` | Yes | — | How often scx-adapt reads system metrics, in milliseconds |
| `path` | Yes | all schedulers | Where the scheduler file is. Can be a relative path, but absolute is safer. Schedulers added with `scheduler add` can be referred to by filename only. |
| `loader` | Yes | all schedulers | How the scheduler is loaded: `external` (raw BPF bytecode, scx-adapt links it) or `builtin` (standalone executable). |
| `structName` | Only `external` | external | The BPF map key that points to the scheduler's methods. |
| `parameters` | No | `builtin` only | Command-line arguments passed to the `builtin` scheduler. |
| `priority` | Yes | all schedulers | Scheduler priority (1-139). Lower number = checked first. The first scheduler whose criteria match is used. |
| `log` | No | `builtin` only | Turn logging on (`true`) or off (`false`). Scheduler output goes to a file in `/var/log/scx-adapt/`. |
| `criterias` | Yes | all schedulers | The conditions a scheduler must meet to be selected. Each one compares a metric using `more_than` or `less_than`. |
| `value_name` | Yes | within criteria | Which metric to check. See the table below. |

### Currently supported metrics (`value_name`)

| Metric | Format | What it measures |
| --- | --- | --- |
| Pressure Stall Information | `(cpu\|io\|mem)_psi_(some\|full)_(10\|60\|300)` | How much time tasks spent waiting on CPU, I/O or memory |
| Load average | `load_avg_(1\|5\|15)` | System load over the last 1, 5 or 15 minutes |
| Running processes | `procs_running` | How many processes are currently running |
| Blocked processes | `procs_blocked` | How many processes are waiting on I/O |
| Disk I/O | `procs_disk_io` | How many processes are doing I/O right now |

---

## Using scx-adapt

- `scx-adapt [command]`

Commands are grouped into three main areas: `profile`, `scheduler` and `service`.

### Profiles

| Command | What it does |
| --- | --- |
| `profile add <path>` | Copy a profile file into the profiles folder |
| `profile ls` | List all valid profiles in the profiles folder |
| `profile check <path>` | Validate a profile file without adding it |
| `profile edit <filename>` | Open a stored profile in your editor |
| `profile start <path>` | Run scx-adapt with a profile (by filename or any path) |
| `profile rm <filename>` | Delete a profile from the profiles folder |

### Schedulers

| Command | What it does |
| --- | --- |
| `scheduler add --loader external\|builtin <path(s)>` | Add one or more schedulers to the schedulers folder |
| `scheduler ls` | List added schedulers (grouped by loader type) |
| `scheduler rm --loader external\|builtin <filename>` | Remove a scheduler |

### Service

| Command | What it does |
| --- | --- |
| `service install` | Install the systemd service file |
| `service edit` | Edit the service file with `systemctl edit --full` |
| `service reset` | Reset the service file to its defaults |
| `service rm` | Remove the systemd service file |

### Everything else

| Command | What it does |
| --- | --- |
| `check-dependencies` | Check that the kernel and BPF setup are ready |
| `status` | Show scheduler/profile counts and current state |
| `trace-sched <file>` | Write sched_ext event tracing to a file |

---

## Systemd service

```ini
[Unit]
Description=scx-adapt daemon for profile at %I
StartLimitIntervalSec=30
StartLimitBurst=4 

[Service]
Type=exec
ExecStart=/usr/bin/scx-adapt profile start %i
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

scx-adapt service self-heals unless schedulers, which are defined inside the YAML configuration file, fail **more than four times within 30 seconds** (by default).  

To install the service file and enable/start the service:
- `# scx-adapt service install`
- `# systemctl enable scx-adapt@<profile_path>`
- `# systemctl start scx-adapt@<profile_path>`

To disable/stop and delete service file:
 
- `# systemctl disable scx-adapt@<profile_path>`
- `# systemctl stop scx-adapt@<profile_path>`
- `# scx-adapt service rm`

---

## Logs

For builtin-loader schedulers with `log: true`, standard output and standard error streams are written into log files which are located in `/var/log/scx-adapt` by default.

```bash
captain@fedora:/var/log/scx-adapt 🚢🐳 $ ls -1
2026-05-10_16-51-23_scx_central.log
2026-05-10_17-27-38_scx_beerland.log
2026-05-10_17-32-55_scx_rustland.log
2026-05-10_17-35-35_scx_simple.log
2026-05-10_17-47-34_scx_chaos.log
```

---

## Installation

### go install

`go install github.com/dogukanmeral/scx-adapt@latest`

> Pre-compiled executables and distribution packages are not available but I plan to add automated releases as soon as possible. 

### Dependencies

The kernel has to be built with the following configuration:
- CONFIG_BPF=y
- CONFIG_BPF_SYSCALL=y
- CONFIG_BPF_JIT=y
- CONFIG_DEBUG_INFO_BTF=y
- CONFIG_BPF_JIT_ALWAYS_ON=y
- CONFIG_BPF_JIT_DEFAULT_ON=y
- CONFIG_SCHED_CLASS_EXT=y

<https://github.com/sched-ext/scx#build--install>

---

## Further development

I (Doğukan Meral) have been the sole developer for scx-adapt while my friend [Onur Karagür](https://github.com/onurkaragur/) helped me on performance analysis of schedulers and ways to optimize scx-adapt configurations using machine learning techniques on the [scx-adapt-experiments](https://github.com/onurkaragur/scx-adapt-experiments) repository.

Your feedbacks, suggestions, criticisms and most importantly your contributions are highly appreciated. Feel free to contact me at my e-mail address `dogukan.meral@protonmail.com`   

---

## Helpful resources for bpf and sched_ext / Inspirations for this project

- Johannes Bechberger's [minimal scheduler repository](https://github.com/parttimenerd/minimal-scheduler) and [blog article](https://mostlynerdless.de/blog/2024/10/25/a-minimal-scheduler-with-ebpf-sched_ext-and-c/)
- Andrea Righi's [neural network scheduler video](https://youtu.be/ywW83YK9EKQ)
- [scx repository](https://github.com/sched-ext/scx) which contains many Sched_ext schedulers and tools
- [Perfetto](https://ui.perfetto.dev/): Browser based and locally running scheduler trace visualisation and analysis tool
- Changwoo Min's ['sched_ext: a BPF-extensible scheduler class'](https://blogs.igalia.com/changwoo/sched-ext-a-bpf-extensible-scheduler-class-part-1/) and ['sched_ext: scheduler architecture and interfaces'](blogs.igalia.com/changwoo/sched-ext-scheduler-architecture-and-interfaces-part-2/) blog articles

---

## License

scx-adapt is licensed under GPLv2. See LICENSE.