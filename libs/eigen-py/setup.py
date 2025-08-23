from setuptools import setup, find_packages

setup(
    name='eigen-py',
    version='0.1.0',
    package_dir={'': 'src'},
    packages=find_packages(where='src'),
    install_requires=[
        'requests==2.29.0',
        'pytest==8.4.1',
        'transformers==4.47.1'
    ],
)
