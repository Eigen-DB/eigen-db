<div align="center">

# faissgo⚡

### Lightweight Go bindings for [Meta’s Faiss](https://engineering.fb.com/2017/03/29/data-infrastructure/faiss-a-library-for-efficient-similarity-search/) 🚀

![Faiss](https://img.shields.io/badge/Faiss-%230467DF.svg?style=for-the-badge&logo=Meta&logoColor=white)
![CUDA](https://img.shields.io/badge/cuda-000000.svg?style=for-the-badge&logo=nVIDIA&logoColor=green)

[![Go Reference](https://pkg.go.dev/badge/github.com/Eigen-DB/eigen-db/libs/faissgo)](https://pkg.go.dev/github.com/Eigen-DB/eigen-db/libs/faissgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/Eigen-DB/eigen-db)](https://goreportcard.com/badge/github.com/Eigen-DB/eigen-db)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

</div>

## ✨ Features

* Easy-to-use bindings for Faiss, enabling fast and efficient similarity search of high-dimensional vectors. 
* Idiomatic Go bindings for Faiss, allowing seamless integration of Faiss' state-of-the-art vector indexing and search capabilities into Go applications. 
* (COMING SOON) Support's GPU acceleration for blazingly fast vector operations. 

## 📦 Installation

### 1. Install system dependencies:

**Ubuntu/Debian**

```sh
apt update
apt install -y libopenblas-dev libgomp1
```

**Fedora**

```sh
dnf install -y openblas-devel libgomp
```

**Arch Linux**

```sh
pacman -S openblas
```

### 2. Install package

```sh
go get github.com/Eigen-DB/eigen-db/libs/faissgo
```
