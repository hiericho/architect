from setuptools import setup, find_packages

setup(
    name="architect-net",
    version="1.1.1",
    packages=find_packages(),
    description="Low-level network fingerprint evasion engine",
    author="Hiericho",
    python_requires=">=3.7",
    install_requires=[],
    include_package_data=True,
    package_data={
        "architect": ["architect.dll", "architect.so", "architect.dylib"],
    },
)
