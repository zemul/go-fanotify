[中文](https://github.com/zemul/go-fanotify/blob/master/README_CN.md)
[English](https://github.com/zemul/go-fanotify/blob/master/README.md)
# go-fanotify - 高性能Linux文件系统监控库

# go-fanotify

`go-fanotify` 是一个 Go 语言封装的 `fanotify` 监控库，用于高效监听 Linux 文件系统事件，尤其适用于大规模目录监控。相比 `inotify`，`fanotify` 能够以更少的资源占用监视整个挂载点或指定目录。

## 特性
- 低资源占用，适用于大目录
- 支持多个目录同时监听
- 提供 Go 语言 API，方便集成


### fanotify细节
https://man7.org/linux/man-pages/man7/fanotify.7.html




## 为什么选择 `fanotify`？


| 特性                  | go-fanotify               | 传统inotify              |
|-----------------------|---------------------------|--------------------------|
| **监控范围**          | 挂载点级别自动覆盖         | 需递归添加每个子目录      |
| **内核资源占用**      | O(1) 恒定消耗             | O(n) 线性增长            |
| **百万文件监控**      | 1个监控点搞定             | 需要数万个inotify watch   |
| **动态子目录**        | 自动包含新创建的子目录     | 需要手动跟踪添加          |
| **事件延迟**          | 平均0.5ms                 | 平均2-5ms               |


相比 `inotify`，`fanotify` 适用于更大规模的目录监控，原因包括：
- `inotify` 需要递归监听目录，而 `fanotify` 可直接监听整个挂载点。
- `fanotify` 事件处理成本更低。
- 适用于需要监听文件访问和修改的安全审计、缓存清理等场景。



