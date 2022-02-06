FROM redpwn/jail:0.1.3

COPY --from=python:3.9 / /srv
COPY flag.txt /srv/
RUN chmod 444 /srv/flag.txt && mv /srv/flag.txt /srv/flag.`tr -dc A-Za-z0-9 < /dev/urandom | head -c 20`.txt
COPY ti1337se.py /srv/app/run
RUN chmod 755 /srv/app/run
