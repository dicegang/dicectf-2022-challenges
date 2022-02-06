FROM redpwn/jail:0.1.3
ENV JAIL_MEM=0
ENV JAIL_TMP_SIZE=4096
COPY --from=debian:11.2-slim / /srv
RUN mkdir /srv/app
COPY bin/snapshot_blob.bin /srv/app
COPY bin/d8 /srv/app/
COPY wrapper /srv/app/run
COPY bin/asfdjklafdsljkadsf4538df.txt /srv/app
