# 《kubernetes源码阅读》

kubernetes由许多组建构成

# 目录

## kubelet
kubelet是每个节点都要安装的重要组建，向上与apiserver交互，向下通过cri与容器运行时交互，完成...职责。
- [ kubelet ](kubelet/README.md)