# Copyright (C) 2020 Intel Corporation
# SPDX-License-Identifier: BSD-3-Clause

FROM alpine:3.16

RUN apk update && apk add --no-cache \
    bash \
    rsync 

RUN mkdir /ace

COPY /bin /ace/bin
COPY /demo_certs /ace/demo_certs
COPY /systemd /ace/systemd
COPY /dxo /ace/dxo
COPY /stacks /stacks

ENV INSTALL_TYPE=demo \
    APP_INSTALL=true \
    SYSTEMD_INSTALL=false

ENTRYPOINT ["/bin/bash"]
CMD ["/ace/bin/install.sh"]
