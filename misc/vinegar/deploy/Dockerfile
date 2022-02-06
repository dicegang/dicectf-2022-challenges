FROM redpwn/jail:0.1.3

COPY --from=python:3.9 / /srv
COPY vinegar.py /srv/app/run
RUN chmod 755 /srv/app/run
