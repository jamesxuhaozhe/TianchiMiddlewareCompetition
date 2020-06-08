FROM centos:centos7

COPY tailSamp /usr/local/src
COPY start.sh /usr/local/src
WORKDIR /usr/local/src
RUN chmod +x /usr/local/src/start.sh
RUN chmod +x /usr/local/src/tailSamp
ENTRYPOINT ["/bin/bash", "/usr/local/src/start.sh"]
