## ACE Installation

`Autonomous Cluster for the Edge` or **ACE**

### Pre-Requisites

- **hardware requirements**

    * x86 Hardware or x86 Virtual Machine
    * At Least 5 GB of Disk Space
	* 4 GB of RAM

- **software requirements**

    * `docker` 18.06.x or greater
    * `docker-compose` v1.23.2 or greater
    * `bash` v4.3.48 or greater

[Install `docker`](https://docs.docker.com/install/)

[Install `docker-compose`](https://docs.docker.com/compose/install/)

```
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

### Clone the repository

Clone this repo at the path `/opt/ace` using your git protocol of choice.

### Build the ACE Service

`ACE` service needs a working `Internet Connection` during build and run.

####  Normal Build
`ACE` consists of microservices which built as `docker` images.
You need `root` privileges to build the docker `images`.

<a name="build_steps"></a>
**Build Steps**
```
# sudo -i
# cd /opt/ace
# ./build.sh
```

####  Build under corporate proxy

you might need to set proxy for `docker` as

1. Set the appropriate proxy in these variables

`http_proxy`, `https_proxy`, `HTTP_PROXY`, `HTTPS_PROXY`, `HTTPPROXY`, `HTTPSPROXY`

```
export http_proxy=<proxy-url>
export https_proxy=<proxy-url>
export HTTP_PROXY=<proxy-url>
export HTTPS_PROXY=<proxy-url>
export HTTPPROXY=<proxy-url>
export HTTPSPROXY=<proxy-url>
```

2. Follow the [Build Steps](#build_steps)

#### Build Failures

In case of build failures due to network connectivity. please follow the build steps again.


### Set up Security

`ACE` needs an initial configuration to function. This configuration includes *provisioning* of the keys.

#### Provisioning Keys

In this version of `ACE` the key provisioning is **manual**.

`ACE` is designed to *automatically discover* the peer nodes which have `ACE` installed and set up the cluster. The cluster thus formed consists of an [docker-swarm](https://docs.docker.com/engine/swarm/) as an orchestrator and a distributed file storage served by [gluster](https://www.gluster.org/).

**Every Node** in a cluster should have **same** SERF `Key` kept in `/etc/ssl/ace/keyring.json`. SERF Key is a **AES-256** symmetric key.

_Contents of  **/etc/ssl/ace/keyring.json**_
```
[
  "HvY8ubRZMgafUOWvrOadwOckVa1wN3QWAo46FVKbVN8="
]
```
Additionally **RSA PKI Keys** need to be generated and kept at `/etc/ssl`.

Follow the instructions to create all keys [here](02_Security.md).

### Start ACE Service

It is adviseable to set up [ssh console](04_Deploy_Workloads.md#set-up-console) or [UI based swam management tools before starting `ACE`. 


- using systemd

*Install ace.service*

```
# cp /opt/ace/systemd/ace.service /etc/systemd/system/
# ln -s /etc/systemd/system/ace.service /etc/systemd/system/multi-user.target.wants/ace.service
```

*Start Service*

```
systemctl start ace

```

- without systemd

```
docker-compose -p ace -f /opt/ace/compose/docker-compose.yml up -d

```

### Using ACE

`ACE` forms a `docker swarm` cluster of the nodes which are in the network.

Service stacks can be deployed on the cluster using any of the following methods.

#### CLI through SSH server

[See Instructions Here.](04_Deploy_Workloads.md#set-up-console)


#### Web consoles

- Portainer

It can be used to set up `remote management` as well as `local management` of swarm cluster through a WEB UI.

[See Instructions Here.](04_Deploy_Workloads.md#set-up-management-tools)

- Swarmpit

It can be used to set up`local management` of swarm cluster through a WEB UI.


[See Instructions Here.](04_Deploy_Workloads.md#set-up-management-tools)

### Stop ACE Service

- systemd


```
systemctl stop ace

```

- without systemd


```
docker-compose -p ace -f /opt/ace/compose/docker-compose.yml down
```

#### In order to factory reset the node in case of a problem.

```
/opt/ace/bin/reset
```
