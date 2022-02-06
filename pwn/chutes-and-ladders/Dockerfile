FROM redpwn/jail@sha256:c86bba3a1da16b181b63f5e500f6955f3f68ec44caf0acfd33771d2b7f634d4b
# ubuntu 20.04
COPY --from=ubuntu@sha256:b3e2e47d016c08b3396b5ebe06ab0b711c34e7f37b98c9d37abe794b71cea0a2 / /srv
COPY bin/chutes /srv/app/run
COPY bin/flag.txt /srv/app
