FROM debian:latest AS builder

SHELL ["/bin/bash", "-c"]

RUN apt update && apt install -y \
    build-essential cmake ninja-build wget git

#RUN apt install -y apt-key

COPY ./lib /app

WORKDIR /app/faiss

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
#RUN make -C build install (not needed since we are copying the files into the next stage)
#RUN sudo cp build/c_api/libfaiss_c.so /usr/lib

FROM golang:1.23

RUN apt update -y

COPY . /app

WORKDIR /app

# install Intel MKL
RUN wget https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN apt-key add GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN sh -c 'echo "deb https://apt.repos.intel.com/oneapi all main" > /etc/apt/sources.list.d/oneAPI.list'
RUN apt update -y
RUN apt install -y intel-oneapi-mkl

COPY --from=builder /app/faiss/build/c_api/libfaiss_c.so /usr/lib
COPY --from=builder /app/faiss/build/faiss/libfaiss.so /usr/lib
COPY --from=builder /app/faiss /usr/include/faiss

RUN go mod tidy

RUN echo "source /opt/intel/oneapi/mkl/latest/env/vars.sh" >> ~/.bashrc

ENTRYPOINT [ "/bin/bash", "-c"]
