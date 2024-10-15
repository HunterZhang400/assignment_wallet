FROM golang AS golang-build-env
ADD . /go/src/assignment_wallet
WORKDIR /go/src/assignment_wallet
ENV GO111MODULE=on
ENV CGO_ENABLED=1
RUN go build -mod=vendor -v -o /user/bin/assignment_wallet main.go
#RUN apt-get update && apt-get -y install  curl telnet wget iputils-ping bash-completion
EXPOSE 80
RUN chmod +X /user/bin/assignment_wallet
CMD ["/user/bin/assignment_wallet"]