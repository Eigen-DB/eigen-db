# this dockerfile contains the steps to create the proper environment for the monorepo

FROM golang:1.23

# installing golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.5

# choose between MKL or openblas
RUN apt update -y
RUN apt install -y libopenblas-dev libgomp1
RUN wget https://apt.repos.intel.com/intel-gpg-keys/GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN apt-key add GPG-PUB-KEY-INTEL-SW-PRODUCTS.PUB
RUN sh -c 'echo "deb https://apt.repos.intel.com/oneapi all main" > /etc/apt/sources.list.d/oneAPI.list'
RUN apt update -y
RUN apt install -y intel-oneapi-mkl
RUN echo "source /opt/intel/oneapi/mkl/latest/env/vars.sh" >> ~/.bashrc
RUN echo "export PATH=/usr/local/go/bin:/go/bin:$PATH" >> ~/.bashrc

# creating user to use in the CI (some tests expect lower-level system access)
RUN useradd -m ci_user -s /bin/bash

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
