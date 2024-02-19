基于eunomia-bpf开发工具进行

## 安装开发环境
下载ecli工具来运行ebpf程序
```shell
$ wget https://aka.pw/bpf-ecli -O ecli && chmod +x ./ecli
$ ./ecli -h
Usage: ecli [--help] [--version] [--json] [--no-cache] url-and-args
```
安装编译工具链把eBPF代码编译成config files（这个文件可以被ecli使用）或者WASM模块
```shell
$ wget https://github.com/eunomia-bpf/eunomia-bpf/releases/latest/download/ecc && chmod +x ./ecc
$ ./ecc -h
eunomia-bpf compiler
Usage: ecc [OPTIONS] <SOURCE_PATH> [EXPORT_EVENT_HEADER]
....
```
也可以用镜像编译 
```shell
$ docker run -it -v `pwd`/:/src/ ghcr.io/eunomia-bpf/ecc-`uname -m`:latest # Compile using docker. `pwd` should contain *.bpf.c files and *.h files.
export PATH=PATH:~/.eunomia/bin
Compiling bpf object...
Packing ebpf object and config into /src/package.json...
```