# pylint: disable=missing-docstring

import os

from test_provider_base import BaseProviderTest, create_provider_test_metaclass
from test_provider_oks_config import OKS_SERVICE_CONFIG

ROOT_DIR = os.path.join(os.path.dirname(__file__), "data", "oks")

OKSMeta = create_provider_test_metaclass(ROOT_DIR)

class TestProviderOKS(BaseProviderTest, metaclass=OKSMeta):
    service_config = OKS_SERVICE_CONFIG
