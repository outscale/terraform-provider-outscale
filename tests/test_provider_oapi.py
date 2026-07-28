# pylint: disable=missing-docstring

import os

from test_provider_base import BaseProviderTest, create_provider_test_metaclass
from test_provider_oapi_config import OAPI_SERVICE_CONFIG

ROOT_DIR = os.path.join(os.path.dirname(__file__), "data", "oapi")

def oapi_env_filter(resource):
    if os.getenv("SKIP_NETS"):
        return resource != "nets"
    if os.getenv("RUN_NETS_ONLY"):
        return resource == "nets"
    return True

OAPIMeta = create_provider_test_metaclass(ROOT_DIR, oapi_env_filter)

class TestProviderOAPI(BaseProviderTest, metaclass=OAPIMeta):
    service_config = OAPI_SERVICE_CONFIG
