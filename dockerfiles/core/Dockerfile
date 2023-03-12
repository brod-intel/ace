# Copyright (C) 2020 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

FROM docker:20.10.23-cli

RUN apk update && apk add --no-cache \
    bash \
    openssl \
    rsync 

RUN mkdir /ace

COPY /bin /ace/bin
COPY /demo_certs /ace/demo_certs
COPY /systemd /ace/systemd
COPY /dxo /ace/dxo
COPY /gluster /ace/gluster

ENV INSTALL_TYPE=demo \
    APP_INSTALL=true \
    SYSTEMD_INSTALL=false

ENTRYPOINT ["/bin/bash"]
CMD ["/ace/bin/install.sh"]
