from setuptools import setup, find_packages

setup(
    name="KnowEase-RAGService",
    version="0.1",
    packages=find_packages(),
    install_requires=[
        "haystack-ai",
        "elasticsearch-haystack",
        "sqlalchemy",
        "pymysql",
        "grpcio",
        "grpcio-tools",
        "pyyaml"
    ],
    entry_points={
        'console_scripts':[
            'KnowEase='
        ]
    }
) 