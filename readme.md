### base service

docker-compose.yml

registry生成htpasswd认证：

docker run --entrypoint htpasswd \
httpd:2 -Bbn 用户名 密码 > 文件
