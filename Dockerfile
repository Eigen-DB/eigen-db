# this Dockerfile defines the base image for the eigendb monorepo

# setting up faissgo deps
FROM debian:latest AS faissgo-builder

SHELL ["/bin/bash", "-c"]

COPY ./libs/faissgo/lib /faissgo

RUN apt update && apt install -y \
    build-essential cmake ninja-build wget git

WORKDIR /faissgo/faiss

RUN rm -rf build

# Compiling Faiss
RUN wget https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN apt-key add GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN sh -c 'echo "deb https://apt.repos.intel.com/oneapi all main" > /etc/apt/sources.list.d/oneAPI.list'
RUN apt update -y
RUN apt install -y intel-oneapi-mkl
RUN apt install -y libgflags-dev
RUN source /opt/intel/oneapi/mkl/latest/env/vars.sh && cmake -B build \
    -DFAISS_ENABLE_GPU=OFF \
    -DFAISS_ENABLE_C_API=ON \
    -DBUILD_SHARED_LIBS=ON \
    -DFAISS_ENABLE_PYTHON=OFF \
    -DBUILD_TESTING=OFF \
    -DFAISS_ENABLE_MKL=ON \
    .
RUN make -C build -j

FROM golang:1.23

RUN go env -w GOFLAGS='-buildvcs=false'

# install Intel MKL (could probably speed up the build process by just copying in /opt/intel from the faissgo-builder stage)
RUN apt update -y
RUN wget https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN apt-key add GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN sh -c 'echo "deb https://apt.repos.intel.com/oneapi all main" > /etc/apt/sources.list.d/oneAPI.list'
RUN apt update -y
RUN apt install -y intel-oneapi-mkl

COPY --from=faissgo-builder /faissgo/faiss/build/c_api/libfaiss_c.so /usr/lib
COPY --from=faissgo-builder /faissgo/faiss/build/faiss/libfaiss.so /usr/lib
COPY --from=faissgo-builder /faissgo/faiss /usr/include/faiss

RUN echo "source /opt/intel/oneapi/mkl/latest/env/vars.sh" >> ~/.bashrc
RUN echo "export PATH=/usr/local/go/bin:/go/bin:$PATH" >> ~/.bashrc

# creating user for CI
RUN useradd -m ci_user -s /bin/bash
RUN su ci_user -c "curl -fsSL https://moonrepo.dev/install/moon.sh | bash && echo 'export PATH=\"\$HOME/.moon/bin:\$PATH\"' >> ~/.bashrc"

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
