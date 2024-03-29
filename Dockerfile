FROM nexus3.o-ran-sc.org:10004/o-ran-sc/bldr-ubuntu20-c-go:1.1.1 as kpimonbuild

ENV PATH $PATH:/usr/local/bin
ENV GOPATH /go
ENV GOBIN /go/bin
ENV RMR_SEED_RT /opt/routes.txt

COPY routes.txt /opt/routes.txt

#ARG RMRVERSION=4.0.2
#ARG RMRLIBURL=https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr_${RMRVERSION}_amd64.deb/download.deb
#ARG RMRDEVURL=https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_${RMRVERSION}_amd64.deb/download.deb
#RUN wget --content-disposition ${RMRLIBURL} && dpkg -i rmr_${RMRVERSION}_amd64.deb
#RUN wget --content-disposition ${RMRDEVURL} && dpkg -i rmr-dev_${RMRVERSION}_amd64.deb
#RUN rm -f rmr_${RMRVERSION}_amd64.deb rmr-dev_${RMRVERSION}_amd64.deb


RUN wget -nv --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_4.0.5_amd64.deb/download.deb
RUN dpkg -i rmr-dev_4.0.5_amd64.deb

# Install RMR client
RUN apt-get update && \
    apt-get -y install gcc

ARG XAPPFRAMEVERSION=v0.4.11
WORKDIR /go/src/gerrit.o-ran-sc.org/r/ric-plt
# RUN git clone "https://gerrit.o-ran-sc.org/r/ric-plt/sdlgo"
RUN git clone -b ${XAPPFRAMEVERSION} "https://gerrit.o-ran-sc.org/r/ric-plt/xapp-frame"
RUN cd xapp-frame && \
   GO111MODULE=on go mod vendor -v && \
    cp -r vendor/* /go/src/ && \
    rm -rf vendor

WORKDIR /go/src/gerrit.o-ran-sc.org/r/scp/ric-app/kpimon
COPY control/ control/
COPY cmd/kpimon.go ./kpimon.go
COPY e2ap/ e2ap/
COPY e2sm/ e2sm/
COPY ./go.mod ./go.mod

# "COMPILING E2AP Wrapper"
# -DASN_EMIT_DEBUG=1
RUN cd e2ap && \
    gcc -c -fPIC -Iheaders/ lib/*.c wrapper.c && \
    gcc *.o -shared -o libe2apwrapper.so && \
    cp libe2apwrapper.so /usr/local/lib/ && \
    mkdir /usr/local/include/e2ap && \
    cp wrapper.h headers/*.h /usr/local/include/e2ap && \
    ldconfig

# "COMPILING E2SM Wrapper"
RUN cd e2sm && \
    gcc -DASN_EMIT_DEBUG=1 -c -fPIC -Iheaders/ lib/*.c wrapper.c && \
    gcc *.o -shared -o libe2smwrapper.so && \
    cp libe2smwrapper.so /usr/local/lib/ && \
    mkdir /usr/local/include/e2sm && \
    cp wrapper.h headers/*.h /usr/local/include/e2sm && \
    ldconfig

WORKDIR /go/src/gerrit.o-ran-sc.org/r/scp/ric-app/kpimon

COPY control/ control/
COPY cmd/kpimon.go ./kpimon.go
COPY ./go.mod ./go.mod


RUN wget -nv --no-check-certificate https://dl.google.com/go/go1.18.linux-amd64.tar.gz \
     && tar -xf go1.18.linux-amd64.tar.gz \
     && rm -f go*.gz

ENV DEFAULTPATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV PATH=$DEFAULTPATH:/usr/local/go/bin:/go/src/gerrit.o-ran-sc.org/r/scp/ric-app/kpimon/go/bin:/root/go/bin

COPY go.sum go.sum

RUN mkdir pkg

#RUN go env -w GO111MODULE=off
RUN go build ./kpimon.go && pwd && ls -lat

FROM ubuntu:20.04
COPY --from=kpimonbuild /usr/local/lib /usr/local/lib
COPY --from=kpimonbuild /usr/local/include/e2ap/*.h /usr/local/include/e2ap/
COPY --from=kpimonbuild /usr/local/include/e2sm/*.h /usr/local/include/e2sm/
RUN ldconfig
WORKDIR /go/src/gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/config/
COPY --from=kpimonbuild /go/src/gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/config/config-file.yaml .
WORKDIR /go/src/gerrit.o-ran-sc.org/r/scp/ric-app/kpimon
COPY --from=kpimonbuild /go/src/gerrit.o-ran-sc.org/r/scp/ric-app/kpimon/kpimon .

ENV  RMR_RTG_SVC="9999" \
     VERBOSE=0 \
     CONFIG_FILE=/go/src/gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/config/config-file.yaml

CMD ./kpimon -f $CONFIG_FILE
