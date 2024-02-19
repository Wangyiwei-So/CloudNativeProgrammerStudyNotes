#ifndef __MEMLEAK_H
#define __MEMLEAK_H

#define ALLOCS_MAX_ENTRIES 1000000
#define COMBINED_ALLOCS_MAX_ENTRIES 10240
#include <vmlinux.h>
struct alloc_info {
    __u64 size;            // Size of allocated memory
    __u64 timestamp_ns;    // Timestamp when allocation occurs, in nanoseconds
    int stack_id;          // Call stack ID when allocation occurs
};

union combined_alloc_info {
    struct {
        __u64 total_size : 40;        // Total size of all unreleased allocations
        __u64 number_of_allocs : 24;   // Total number of unreleased allocations
    };
    __u64 bits;    // Bitwise representation of the structure
};

#endif /* __MEMLEAK_H */