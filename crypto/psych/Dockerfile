FROM python:slim AS app
RUN pip install --no-cache-dir sibc==1.0.3

FROM redpwn/jail
COPY --from=app / /srv
COPY server.py /srv/app/run
COPY pk.bin sk.bin /srv/app/
ENV JAIL_TIME=60 JAIL_MEM=100M JAIL_CPU=500