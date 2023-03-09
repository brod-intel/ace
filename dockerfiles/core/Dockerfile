# Copyright (C) 2020 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

FROM alpine:3.16

RUN apk update && apk add --no-cache \
    bash \
    rsync 

RUN mkdir /opt/ace

COPY /bin /opt/ace/bin
COPY /demo_certs /opt/ace/demo_certs
COPY /systemd /opt/ace/systemd
COPY /dxo /opt/ace/dxo

ENV INSTALL_TYPE=demo \
    APP_INSTALL=true \
    SYSTEMD_INSTALL=false

ENTRYPOINT ["/bin/bash"]
CMD ["/opt/ace/bin/install.sh"]
