FROM python:slim AS app

FROM redpwn/jail
COPY --from=app / /srv
COPY --chmod=755 correlated.py /srv/app/run
COPY flag.txt /srv/app/
ENV JAIL_MEM=100M JAIL_TIME=60 JAIL_CPU=500
