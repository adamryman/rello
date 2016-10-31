FROM golang:1.7.1

ENV TINI_VERSION v0.10.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
# or docker run your-image /your/program ...

RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64
RUN chmod +x /usr/local/bin/dumb-init

RUN mkdir -p /go/src/github.com/adamryman/rello

COPY . /go/src/github.com/adamryman/rello

WORKDIR /go/src/github.com/adamryman/rello
	
RUN go install -v ./...

EXPOSE 5040

ENV PORT 0
ENV HTTPPORT 5040

#ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]
ENTRYPOINT ["/tini", "--"]
CMD ["/go/bin/rello-server", "-http.addr", ":5040", "-grpc.addr", ":0"]
