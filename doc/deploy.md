# 部署

本项目已经打包成[docker镜像](https://hub.docker.com/r/longfangsong/shu-course-proxy)。

## 支持服务
### postgresql数据库
migration文件位于本repo的 [migration](https://github.com/longfangsong/SHUCourseProxy/tree/master/migration) 目录中。

建议使用 [golang-migrate](https://github.com/golang-migrate/migrate) 来进行 migrate。
```shell
migrate -source github://[你的Github用户名]:[你的Github Access Token]@longfangsong/SHUCourseProxy/migration -database [你的postgrsql数据库url] up
```
## 服务本身
### 环境变量
- `PORT`: 服务端口
- `DB_ADDRESS`: 数据库url
- `JWT_SECRET`: jwt密钥
### k8s
在k8s下使用如下yaml，假设`JWT_SECRET`由k8s secret给出，数据库服务器在`shu-course-proxy-svc`。
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shu-course-proxy
spec:
  selector:
    matchLabels:
      run: shu-course-proxy
  replicas: 1
  template:
    metadata:
      labels:
        run: shu-course-proxy
    spec:
      containers:
      - name: shu-course-proxy
        image: longfangsong/shu-course-proxy
        env:
        - name: PORT
          value: "8086"
        - name: DB_ADDRESS
          value: "postgres://shuosc@shu-course-proxy-postgres-svc:5432/shu-course-proxy?sslmode=disable"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: shuosc-secret
              key: JWT_SECRET
        ports:
        - containerPort: 8086
---
apiVersion: v1
kind: Service
metadata:
  name: shu-course-proxy-svc
spec:
  selector:
    run: shu-course-proxy
  ports:
  - protocol: TCP
    port: 8086
    targetPort: 8086