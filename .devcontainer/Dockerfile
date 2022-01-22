FROM golang:1.17.6-alpine

WORKDIR /workspace

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
    && apk add git zsh curl wget nodejs docker
RUN git config --global url."https://gitclone.com/".insteadOf https://
RUN sh -c "$(curl -fsSL https://gitee.com/mirrors/oh-my-zsh/raw/master/tools/install.sh)"

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go get -u github.com/spf13/cobra/cobra \
    && go get -v github.com/ramya-rao-a/go-outline \
    && go get -v golang.org/x/tools/gopls

ENV SHELL /bin/zsh

CMD ["tail", "-f", "/dev/null"]