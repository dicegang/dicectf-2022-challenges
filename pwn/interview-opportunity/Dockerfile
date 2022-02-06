FROM redpwn/jail:0.1.3
COPY --from=debian:11.2-slim / /srv
COPY bin/interview-opportunity /srv/app/run
COPY bin/flag.txt /srv/app
