FROM redpwn/jail:0.1.3
COPY --from=ubuntu:20.04 / /srv
COPY run.sh /srv/app/run
COPY breach breach.bin flag.txt /srv/app/
