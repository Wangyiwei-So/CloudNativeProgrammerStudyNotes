#include <vmlinux.h>
#include <bpf/bpf_helpers.h>

const volatile int pid_filter = 0;

SEC("tracepoint/syscalls/sys_enter_openat")
int tracepoint_syscalls_sys_enter_openat(struct trace_event_raw_sys_enter* ctx){
    u64 id = bpf_get_current_pid_tgid();
    u32 pid = id >> 32;
    if(pid_filter && pid_filter!=pid)
        return false;
    bpf_printk("Process ID: %d enter sys openat\n",pid);
    return 0;
}

char LICENSE[] SEC("license") = "GPL";
