# 得到最新的 golang docker 镜像
#FROM golang:latest
FROM xcgd/libreoffice
# 在容器内部创建一个目录来存储我们的 web 应用，接着使它成为工作目录。
#RUN mkdir -p /go/src/web-app
WORKDIR /src/b2b

# 复制 web-app 目录到容器中
ADD b2b.tar.gz /src/b2b/
RUN chmod +x /src/b2b/b2b \
&&rm /src/b2b/go.mod \
&&rm /src/b2b/go.sum \
&&rm /src/b2b/pkg -rf \
&&rm /src/b2b/Dockerfile


#COPY b2b /data/b2b
#COPY docs /data/docs



#调整容器时间
#RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#RUN echo "Asia/Shanghai" > /etc/timezone
#RUN dpkg-reconfigure -f noninteractive tzdata


# 下载并安装第三方依赖到容器中
# RUN go get github.com/go-sql-driver/mysql && go get github.com/astaxie/beego && go get github.com/astaxie/beego/logs && go get gopkg.in/gomail.v2 && go get github.com/tealeg/xlsx && go get  github.com/shopspring/decimal
# RUN go build .
# 告诉 Docker 启动容器运行的命令
ENTRYPOINT  ["/src/b2b/b2b"]
