# -*- coding:utf-8 -*-
# pylint: disable=missing-docstring

OKS_IGNORE_END_ELEMENTS = [
    "request_id",
    "created_at",
    "updated_at",
    "kubeconfig",
    "expiration_date",
]

OKS_IGNORE_TYPE_ELEMENTS = {
    "outscale_oks_project": "name",
    "outscale_oks_cluster": "name",
}

OKS_ID_PREFIXES = []

OKS_SERVICE_CONFIG = {
    "ignore_end_elements": OKS_IGNORE_END_ELEMENTS,
    "ignore_type_elements": OKS_IGNORE_TYPE_ELEMENTS,
    "id_prefixes": OKS_ID_PREFIXES,
}
