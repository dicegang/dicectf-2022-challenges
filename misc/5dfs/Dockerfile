FROM ubuntu:focal
RUN apt update && apt install -y python3 fuse
RUN useradd -ms /bin/bash ctf
WORKDIR /home/ctf
COPY ./5dfswmtt ./setup.sh ./
COPY ./flag /
CMD ["/home/ctf/setup.sh"]
