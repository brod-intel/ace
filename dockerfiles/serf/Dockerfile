FROM alpine:3.12 as hashibase
LABEL maintainer Bryan J Rodriguez <bryan.j.rodriguez@intel.com>

RUN set -x \
 && apk add --no-cache \
   git \
   go \ 
   gnupg \
   libc-dev \
 && gpg --version

WORKDIR /tmp

RUN git clone https://github.com/hashicorp/serf.git \
    && cd serf/ && git checkout v0.8.4 \
    && cd cmd/serf && go build \
    && cp serf /usr/local/bin

FROM alpine:3.12

ARG SERF_GID=7946
ARG SERF_UID=7373

RUN set -x \
 && apk add --no-cache \
    coreutils \
    dumb-init \
    jq \
    su-exec \
 && addgroup -g ${SERF_GID} serf \
 && adduser -S -G serf -u ${SERF_UID} serf

COPY --from=hashibase /usr/local/bin/* /usr/local/bin/
COPY serf/ /opt/ace/serf/
COPY gluster/ /opt/ace/gluster/
COPY bin/ /opt/ace/bin/

WORKDIR /

ENTRYPOINT ["dumb-init", "--", "su-exec", "serf:serf", "serf"]
CMD ["help"]
