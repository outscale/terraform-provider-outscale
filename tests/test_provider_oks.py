# -*- coding:utf-8 -*-
# pylint: disable=missing-docstring

import os
from test_provider_base import create_provider_test_metaclass, BaseProviderTest
from test_provider_oks_config import OKS_SERVICE_CONFIG

ROOT_DIR = os.path.join(os.path.dirname(__file__), "data", "oks")

OksMeta = create_provider_test_metaclass(ROOT_DIR)

class TestProviderOks(BaseProviderTest, metaclass=OksMeta):
    service_config = OKS_SERVICE_CONFIG
