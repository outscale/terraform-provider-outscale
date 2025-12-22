#!/usr/bin/python

import json
import logging
import pprint
import sys


def looking_for_attributes(resource):
    """Fetch attributes dict in resource
        In the following example resources.instances[0].attributes will be returned
    {
      "version": 4,
      "terraform_version": "0.12.16",
      "serial": 3,
      "lineage": "5e62916e-5910-65a9-a2a5-4463f0c9a308",
      "outputs": {},
      "resources": [
        {
          "mode": "managed",
          "type": "outscale_vm",
          "name": "vm001",
          "provider": "provider.outscale",
          "instances": [
            {
              "schema_version": 0,
              "attributes": {
                "admin_password": "",
                "architecture": "x86_64",
                "block_device_mappings": [],
        ....


    """
    target_resource_key = [
        res
        for res in resource.keys()
        if res not in ["mode", "type", "name", "provider", "each"]
    ]
    if len(target_resource_key) != 1:
        print("Can not guess resource in {}".format(resource.keys()))
        return None
    target_resource_key = target_resource_key[0]
    if not isinstance(resource[target_resource_key], list):
        print("Expecting resource to be a list")
        return None
    target_resource = resource[target_resource_key][0]
    attributes = target_resource["attributes"]
    return attributes


def flattern(content, line, result):
    if isinstance(content, (bool, int, str)):
        result.append("{}.{}".format(line, content))  # end
        return
    for left, right in content.items():
        if isinstance(right, (bool, int, str)) or right is None:
            result.append("{}.{}".format(line, left))  # end
        elif isinstance(right, list) and len(right) != 0:
            result.append("{}.{}.#".format(line, left))
            for rright_count, rright in enumerate(right):
                current_line = "{}.{}.{}".format(line, left, rright_count)
                flattern(rright, current_line, result)
        elif isinstance(right, list) and len(right) == 0:
            result.append("{}.{}".format(line, left))


def parse_terraform_state_pull(reported, item):
    """
    :param str reported: filepath to terraform state pull file.
    :param str item: item you want to check
    """
    result = []
    with open(reported, "r") as report:
        reported_content = json.loads(report.read())
    resource_item_content = [
        resource
        for resource in reported_content["resources"]
        if resource["type"] == item
    ]
    if len(resource_item_content) == 0:
        print("ERROR. Looking for type={}. See content:\n".format(item))
        pprint.pprint(reported_content)
        return 1
    resource_item_content = resource_item_content[0]
    attributes = looking_for_attributes(resource_item_content)
    flattern(attributes, "", result)
    result = [res[1:] for res in result]
    result.sort()
    return result


def patch_item(item):
    """ugly way to patch without modifying all existing runCmd(./check.py ... )
    Terraform0.11 has items like outscale_public_ip.outscale_public_ip where terraform 0.12 "type" field is just outscale_public_ip
    """
    items = item.split(".")
    if len(items) == 2:
        return items[1]
    return item


def main(reported, attended, item):
    """
    :param str reported: filepath to terraform state pull file.
    :param str attended: filepath to file which lists datasources attributes
    :param str item: item you want to check
    """
    logger = logging.getLogger("tpd_test")
    item = patch_item(item)
    reported_result = parse_terraform_state_pull(reported, item)
    attended_result = []
    with open(attended, "r") as att:
        attended_result = [
            line for line in att.read().splitlines() if not line.startswith("#")
        ]
    if reported_result == 1:
        logger.debug("???")
        return ["unknown error ???"]
    missing = set(attended_result) - set(reported_result)
    unknown = set(reported_result) - set(attended_result)
    if unknown or missing:
        logger.debug(
            "===== check.py : Differences between terraform reported (pull result) attributes and attended attributes"
        )
    for item in unknown:
        logger.debug("Unknown: {}".format(item))
    for item in missing:
        logger.debug("Missing: {}".format(item))
    if len(missing) != 0:
        return missing
    return None


def check():
    """Ce patch degueu pour eviter de commenter tous les tests avec from check import check"""
    pass


if __name__ == "__main__":
    if len(sys.argv) == 4:
        return_code = main(sys.argv[1], sys.argv[2], sys.argv[3])
        sys.exit(return_code)
    else:
        print("command: check.py reportedAttributes attendedAttributes item")
