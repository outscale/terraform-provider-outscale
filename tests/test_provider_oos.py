# pylint: disable=missing-docstring

import os

from test_provider_base import BaseProviderTest, create_provider_test_metaclass
from test_provider_oos_config import OOS_SERVICE_CONFIG

ROOT_DIR = os.path.join(os.path.dirname(__file__), "data", "oos")

OOSMeta = create_provider_test_metaclass(ROOT_DIR)

class TestProviderOOS(BaseProviderTest, metaclass=OOSMeta):
    service_config = OOS_SERVICE_CONFIG
