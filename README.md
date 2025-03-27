# go-fanotify - High-Performance Linux Filesystem Monitoring Library

## go-fanotify

`go-fanotify` is a Go language wrapper for the `fanotify` monitoring library, designed for efficient Linux filesystem event listening. It is particularly useful for large-scale directory monitoring. Compared to `inotify`, `fanotify` can monitor an entire mount point or a specified directory with lower resource consumption.

## Features
- Low resource usage, suitable for large directories
- Supports monitoring multiple directories simultaneously
- Provides a Go API for easy integration

### fanotify Details
[fanotify man page](https://man7.org/linux/man-pages/man7/fanotify.7.html)

## Why Choose `fanotify`?

| Feature               | go-fanotify               | Traditional inotify      |  
|-----------------------|--------------------------|--------------------------|  
| **Monitoring Scope**  | Covers entire mount point | Requires recursive addition of subdirectories |  
| **Kernel Resource Usage** | O(1) constant consumption | O(n) linear growth |  
| **Monitoring Millions of Files** | Single monitor point | Requires thousands of inotify watches |  
| **Dynamic Subdirectories** | Automatically includes newly created subdirectories | Must be manually tracked and added |  
| **Event Latency** | ~0.5ms on average | 2-5ms on average |  

### Advantages Over `inotify`
Compared to `inotify`, `fanotify` is more suitable for large-scale directory monitoring due to:
- `inotify` requiring recursive directory monitoring, while `fanotify` can directly watch an entire mount point.
- Lower event processing overhead with `fanotify`.
- Ideal for use cases such as security auditing, cache invalidation, and file access monitoring.  
