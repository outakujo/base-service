### base service

docker-compose.yml

registry生成htpasswd认证：

docker run --entrypoint htpasswd \
httpd:2 -Bbn 用户名 密码 > 文件

elasticsearch 启动失败权限解决办法：

进入docker的虚拟机，比如rancher或者wsl的虚拟机

```shell
useradd elasticsearch
chown -R elasticsearch:elasticsearch 虚拟机内的挂载目录
```
