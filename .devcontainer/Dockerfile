FROM golang:1.22.0

RUN apt update -y
RUN apt upgrade -y
RUN apt install -y sudo curl git make vim

RUN groupadd -g 1000 devcontainer
RUN useradd -s /bin/bash -m -u 1000 -g 1000 devcontainer
USER devcontainer
