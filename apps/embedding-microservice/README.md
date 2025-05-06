<div align="center">

# ✨ Embedding Microservice 🧠

### 🤗 Powered by Hugging Face 🤗

[![Hugging Face](https://img.shields.io/badge/Hugging%20Face-FFD21E?logo=huggingface&logoColor=000)](#)
[![NumPy](https://img.shields.io/badge/NumPy-4DABCF?logo=numpy&logoColor=fff)](#)
[![FastAPI](https://img.shields.io/badge/FastAPI-009485.svg?logo=fastapi&logoColor=white)](#)

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

</div>

## Overview

The Embedding Microservice is a service that sits on-top of EigenDB, in charge of embedding various forms of data using pretrained machine learning models.

Model inference is done remotely using Hugging Face's [Inference API](https://huggingface.co/docs/huggingface_hub/en/guides/inference) 🤗.

## Purpose

The purpose of this addition to EigenDB is to enable developers to achieve the benefits of similarity search while not requiring any AI/ML knowledge or experience. Using EigenDB without the Embedding Microservice requires developers to insert embeddings themselves via EigenDB's REST API. This forces the developers to handle the process using machine learning to generate embeddings using their data.

Check out the [EigenDB docs](https://eigendb.mintlify.app/) for more details!

---

Made with ❤️ by developers, for developers.