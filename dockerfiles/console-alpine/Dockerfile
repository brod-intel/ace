FROM alpine:3.16
LABEL maintainer Bryan J Rodriguez <bryan.j.rodriguez@intel.com>

RUN set -x \
	&& apk add --no-cache \
	   coreutils \
	   openssh \
	   sudo \
	   dumb-init \
	   bash \
	   curl \
	   rsync \
	   util-linux \
	&& mkdir -p /var/run/sshd \
	&& mkdir -p /etc/ssh
	
ENV PATH /usr/local/bin:/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin
ENV USERNAME ""
ENV PASSWORD ""
ENV AUTHORIZED_KEYS ""

COPY init/alpine-console-init.sh /bin/init
	
WORKDIR /

ENTRYPOINT ["dumb-init", "/bin/init"]
CMD []