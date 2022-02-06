FROM redpwn/jail:0.1.3

COPY --from=node:17-buster-slim / /srv
COPY --chmod=444 flag.txt /srv/
COPY --chmod=755 index.js /srv/app/run

ENV JAIL_PIDS=20
ENV JAIL_MEM=20M
