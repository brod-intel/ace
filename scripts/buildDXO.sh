mkdir -P ./work
cd work/
git clone --depth=1 https://github.com/brod-intel/dxo.git dxo && \
cd dxo && \
docker build ${DOCKER_BUILD_ARGS} -t ace/dxo:1.0 -f build/Dockerfile .