from setuptools import setup, find_packages

setup(
    name='eigen-client',
    version='1.0.2',
    package_dir={'': 'src'},
    packages=find_packages(where='src'),
    install_requires=[
        "requests==2.29.0",
        "pytest==8.4.1",
        "tiktoken==0.11.0",
        "ollama==0.5.3",
        "openai==1.101.0"
    ],
)
