FROM debian:unstable-20220125 as build

RUN echo 'deb http://deb.debian.org/debian experimental main' > /etc/apt/sources.list.d/experimental.list && \
    apt-get -y update && \
    apt-get install -y libc-bin=2.34-0experimental2 libc6=2.34-0experimental2 \
      libc-dev-bin=2.34-0experimental2 libc6-dev=2.34-0experimental2 \
      gcc make build-essential strace && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY bin/* /app/
RUN gcc -fno-stack-protector -Wl,-z,relro,-z,now -no-pie ./uaf.c -o run && \
    strace ./run

FROM debian:unstable-20220125 as app

RUN echo 'deb http://deb.debian.org/debian experimental main' > /etc/apt/sources.list.d/experimental.list && \
    apt-get -y update && \
    apt-get install -y libc-bin=2.34-0experimental2 libc6=2.34-0experimental2 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=build /app/flag.txt /app/run /app/

FROM redpwn/jail
COPY --from=app / /srv
