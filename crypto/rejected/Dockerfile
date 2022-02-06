FROM python:slim AS app

FROM redpwn/jail
COPY --from=app / /srv
COPY --chmod=755 rejected.py /srv/app/run
COPY flag.txt /srv/app/
ENV JAIL_MEM=20M JAIL_TIME=500
