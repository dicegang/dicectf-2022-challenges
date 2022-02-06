FROM python:3.10-buster AS build
RUN pip wheel git+https://github.com/sarojaerabelli/py-fhe@c84d3e9fe4c31d419c904cfd8b8c5fd6d9d8d19 --wheel-dir=/wheels

FROM python:3.10-slim-buster AS app
COPY --from=build /wheels /wheels
RUN pip install /wheels/* && rm -rf /wheels

FROM redpwn/jail:0.1.3

COPY --from=app / /srv
COPY challenge.py flag.txt /srv/app/
COPY --chmod=755 server.py /srv/app/run

ENV JAIL_TIME=120 JAIL_MEM=100M JAIL_CPU=500
