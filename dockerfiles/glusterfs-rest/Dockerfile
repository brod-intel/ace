FROM ace/glusterfs:7
LABEL maintainer Bryan J Rodriguez <bryan.j.rodriguez@intel.com>

ARG GLUSTERREST_PASSWORD="password"

RUN yum --enablerepo=extras install -y epel-release \
	&& yum install -y \
	   python36-pip \
	&& yum clean all \
	&& rm -rf /var/cache/yum \
	&& pip3.6 install --upgrade pip==21.0.1
RUN pip3.6 install defusedxml==0.6.0
	
COPY glusterfs-rest/ /tmp/glusterfs-rest

WORKDIR /tmp/glusterfs-rest/

RUN python3.6 setup.py install \
	&& glusterrest install; true \
	&& rm -fr /tmp/glusterfs-rest/ \
	&& ln -s /usr/local/bin/gunicorn /usr/bin/gunicorn

ENV PORT 5000

COPY init/glusterfs-rest-init.sh /bin/init

WORKDIR /

ENTRYPOINT ["dumb-init", "/bin/init"]
CMD []
