import os
import sys

from setuptools import find_packages, setup

sys.path.insert(0, os.path.abspath('.'))  # isort:skip
from qa_provider_oapi import version  # isort:skip


def parse_requirements(filename):
    """ load requirements from a pip requirements file """
    lineiter = (line.strip() for line in open(filename))
    return [line for line in lineiter if line and not line.startswith("#")]


VERSION = version.__version__
NAME = "osc-qa-provider-oapi"
if version.__branch__:
    NAME = "{}-{}".format(NAME, version.__branch__)

INSTALL_REQUIRES = parse_requirements('requirements.txt')


DEPS = {}
PACKAGES = []
PKG_DIR = {}
PACKAGE_DATA = {}

for DEP in DEPS:
    print('dep = {}'.format(DEP))
    if DEPS[DEP]:
        for SUB_DEP in DEPS[DEP]:
            PKG_DIR[SUB_DEP] = './{}/{}'.format(DEP, SUB_DEP)
        PACKAGES += [p for p in find_packages(where='./{}'.format(DEP)) if p in DEPS[DEP] or p.split('.')[0] in DEPS[DEP]]
    else:
        PKG_DIR[DEP] = './{}/{}'.format(DEP, DEP)
        tmp_pkgs = find_packages(where='./{}'.format(DEP))
        if tmp_pkgs:
            PACKAGES += [p for p in find_packages(where='./{}'.format(DEP)) if p == DEP or p.startswith('{}.'.format(DEP))]
        else:
            PACKAGES.append(DEP)
PACKAGES += find_packages(exclude=['tests'])

setup(
    name=NAME,
    version=VERSION,
    url='http://www.outscale.com',
    author="3DS Outscale QA Team",
    author_email="qa@outscale.com",
    description="3DS Outscale Terraform tests",
    packages=PACKAGES,
    package_dir=PKG_DIR,
    package_data={'qa_provider_oapi': ['data/*/*/*.tf', 'data/*/*/*.ref',
                                       'data/*/*/*.check', 'data/*/*/*.txt']},
    data_files=[('etc/osc-qa-provider-oapi', ['qa_provider_oapi/provider.auto.tfvars',
                                              'qa_provider_oapi/variables.tf',
                                              'qa_provider_oapi/provider.tf',
                                              'qa_provider_oapi/resources.auto.tfvars'])],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
    #scripts=['qa_tina_redwires/bin/osc-qa-tina-redwires'],
    install_requires=INSTALL_REQUIRES
)
