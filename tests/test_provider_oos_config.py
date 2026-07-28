# pylint: disable=missing-docstring

OOS_IGNORE_END_ELEMENTS = [
    "creation_date",
    "display_name",
    "email_address",
    "email_addresses",
    "ids"
]

OOS_IGNORE_TYPE_ELEMENTS = {}

OOS_ID_ATTRIBUTE_PATHS = [
    "bucket",
    "name",
]

OOS_SERVICE_CONFIG = {
    "ignore_end_elements": OOS_IGNORE_END_ELEMENTS,
    "ignore_type_elements": OOS_IGNORE_TYPE_ELEMENTS,
    "id_attribute_paths": OOS_ID_ATTRIBUTE_PATHS,
}
