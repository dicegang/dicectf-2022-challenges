# syntax=docker/dockerfile:1
FROM gcc:9.3
WORKDIR /attack
COPY attack .
RUN gcc -O0 -o run_attack attack.c
RUN adduser limiteduser --shell=/bin/false --system --disabled-login --no-create-home
USER limiteduser
CMD ["./run_attack"]
