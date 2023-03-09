FROM alpine:3.12
LABEL maintainer Bryan J Rodriguez <bryan.j.rodriguez@intel.com>

RUN set -x \
	&& apk add --no-cache \
		bash \
		coreutils \
		dmidecode \
		dumb-init \
		lshw \
		pciutils \
		util-linux

COPY bin/ /opt/ace/dho/
	
WORKDIR /

ENTRYPOINT ["dumb-init"]
CMD ["/opt/ace/dho/init.sh"]
