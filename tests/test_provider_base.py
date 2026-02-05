# -*- coding:utf-8 -*-
# pylint: disable=missing-docstring

import json
import logging
import os
import shutil
import subprocess

# Generic Terraform-related ignore patterns
IGNORE_PATHS = [
    ".lineage",
    ".serial",
    ".terraform_version",
    ".version",
    ".resources.provider",
    ".check_results",
]
TAG_END_PATHS = [".id", "_id"]
NO_TEST_VALUE = "########"
IGNORE_RESOURCE_TYPES = [
    "random_string",
    "random_integer",
]
IGNORE_END_PATHS = []
VARIABLES_FILE_NAME = ["resources.auto.tfvars"]
VARIABLES = ["region"]
SET_KEY_VALUES = ["resources", "tags"]
ID_PREFIX = "##id-"
ID_SUFFIX = "##"

LOG_HANDLER = logging.StreamHandler()
FORMATTER = logging.Formatter(
    "[%(asctime)s] "
    + "[%(levelname)8s]"
    + "[%(module)s.%(funcName)s():%(lineno)d]: "
    + "%(message)s",
    "%m/%d/%Y %H:%M:%S",
)
LOG_HANDLER.setFormatter(FORMATTER)
logging.basicConfig(level=logging.DEBUG, handlers=[LOG_HANDLER])
logging.getLogger("tpd_test").setLevel(logging.DEBUG)

terraform_vars = {}
for file_name in VARIABLES_FILE_NAME:
    file_path = os.path.abspath(os.path.join(os.path.dirname(__file__), file_name))
    if os.path.exists(file_path):
        with open(file_path, "r") as var_file:
            lines = var_file.readlines()
            for line in lines:
                line = line.strip()
                if line.startswith("#"):
                    continue
                elts = line.split("=")
                if len(elts) != 2:
                    continue
                terraform_vars[elts[0].strip()] = elts[1].strip().strip('"')


def get_test_file_names(test_path, prefix="step", suffix=".tf"):
    ret_file_names = []
    for tmp_file in os.listdir(test_path):
        if tmp_file.startswith(prefix) and tmp_file.endswith(suffix):
            ret_file_names.append(tmp_file)
    return sorted(ret_file_names)


def validate_value_ref(path, parent, value, ids, service_config):
    replace_value = None
    replace = None
    id_prefixes = service_config.get("id_prefixes", [])

    if type(value) is str:
        value_items = value.split("-")
        if (
            len(value_items) == 2
            and len(value_items[1]) in [8, 32]
            and value_items[0] in id_prefixes
        ):
            replace = value

    for p in TAG_END_PATHS:
        if path.endswith(p):
            replace = value
            break
    if replace:
        if replace not in ids:
            ids[replace] = "{}{}{}".format(ID_PREFIX, len(ids.keys()), ID_SUFFIX)
        replace_value = ids[replace]

    if not replace_value and type(value) is str:
        tmp_value = value
        for var in terraform_vars:
            if terraform_vars[var] in tmp_value:
                tmp_value = tmp_value.replace(
                    terraform_vars[var], "###{}###".format(var)
                )
        if tmp_value != value:
            replace_value = tmp_value

    if replace_value:
        path_end = path.split(".")[-1]
        if type(parent) is list:
            parent[int(path_end)] = replace_value
        else:
            parent[path_end] = replace_value


def validate_dict_ref(path, json_ref, ids, service_config):
    for key, value in json_ref.items():
        validate_ref("{}.{}".format(path, key), json_ref, value, ids, service_config)


def validate_list_ref(path, json_ref, ids, service_config):
    for i in range(len(json_ref)):
        validate_ref(
            "{}.{}".format(path, i), json_ref, json_ref[i], ids, service_config
        )


def validate_ref(path, parent, value, ids, service_config):
    ignore_end_elements = service_config.get("ignore_end_elements", [])
    ignore_type_elements = service_config.get("ignore_type_elements", {})

    path_end = path.split(".")[-1]
    if path in IGNORE_PATHS or path_end in ignore_end_elements:
        if parent is not None:
            parent[path_end] = NO_TEST_VALUE
        return None
    for p in IGNORE_END_PATHS:
        if path.endswith(p):
            if parent is not None:
                parent[path_end] = NO_TEST_VALUE
            return None

    if type(value) is list and path == ".resources":
        # Remove resources that we do not want to store in the ref file
        value[:] = [
            resource
            for resource in value
            if not (
                isinstance(resource, dict)
                and resource.get("type") in IGNORE_RESOURCE_TYPES
            )
        ]

    if path_end == "type" and value in ignore_type_elements:
        ignored_key = ignore_type_elements[value]

        if parent and "instances" in parent:
            for instance in parent["instances"]:
                if "attributes" in instance and ignored_key in instance["attributes"]:
                    instance["attributes"][ignored_key] = NO_TEST_VALUE

    if type(value) is dict:
        validate_dict_ref(path, value, ids, service_config)
    elif type(value) is list:
        validate_list_ref(path, value, ids, service_config)
    elif type(value) is tuple:
        assert False, "Unexpected type tuple for path {}".format(path)
    else:
        validate_value_ref(path, parent, value, ids, service_config)
    return value


def compare_json_dicts(path, dict_out, dict_ref, ids):
    keys_out = sorted(set(dict_out.keys()))
    keys_ref = sorted(set(dict_ref.keys()))
    assert len(keys_out) == len(keys_ref), (
        "Not the same keys number for path {}".format(path)
    )
    for key in keys_out:
        assert key in keys_ref, "Could not find key {}.{} in output".format(path, key)
    for key in keys_ref:
        assert key in keys_out, "Could not find key {}.{} in reference".format(
            path, key
        )
    for key in keys_out:
        do_set = False
        if key in SET_KEY_VALUES:
            do_set = True
        elif (
            key.endswith("s")
            and type(dict_out[key]) is list
            and dict_out[key]
            and type(dict_out[key][0]) is dict
            and "{}_id".format(key[:-1]) in dict_out[key][0]
        ):
            do_set = True
        if do_set:
            compare_json_sets(
                "{}.{}".format(path, key), dict_out[key], dict_ref[key], ids
            )
        else:
            compare_json("{}.{}".format(path, key), dict_out[key], dict_ref[key], ids)


def compare_json_lists(path, list_out, list_ref, ids):
    assert len(list_out) == len(list_ref)
    current_ids = ids.copy()
    found_elts = []
    for out_elt in list_out:
        errors = []
        for ref_elt in list_ref:
            if ref_elt in found_elts:
                continue
            try:
                tmp_ids = current_ids.copy()
                compare_json("{}".format(path), out_elt, ref_elt, tmp_ids)
                found_elts.append(ref_elt)
                current_ids = tmp_ids
                errors = []
                break
            except AssertionError as error:
                errors.append(error)
        if errors:
            assert False, "Could not match list values for path {}, {}".format(
                path, errors
            )


def compare_json_sets(path, set_out, set_ref, ids):
    assert len(set_out) == len(set_ref)
    current_ids = ids.copy()
    found_elts = []
    for out_elt in set_out:
        errors = []
        for ref_elt in set_ref:
            if ref_elt in found_elts:
                continue
            try:
                tmp_ids = current_ids.copy()
                compare_json("{}".format(path), out_elt, ref_elt, tmp_ids)
                found_elts.append(ref_elt)
                current_ids = tmp_ids
                errors = []
                break
            except AssertionError as error:
                errors.append(error)
        if errors:
            assert False, "Could not match set values for path {}, {}".format(
                path, errors
            )


def compare_json_values(path, val_out, val_ref, ids):
    # do not check values that should be ignored
    if val_ref == NO_TEST_VALUE:
        return

    if val_out == val_ref:
        return

    # accept id change
    if val_out in ids and ids[val_out] == val_ref:
        return

    if (
        val_out not in ids
        and val_out.startswith(ID_PREFIX)
        and val_ref.startswith(ID_PREFIX)
    ):
        ids[val_out] = val_ref
        return

    # accept string values that only differ by the last digit
    if (
        type(val_out) is str
        and len(val_out) > 1
        and len(val_ref) > 1
        and val_out[:-1] == val_ref[:-1]
    ):
        try:
            int(val_out[-1])
            int(val_ref[-1])
            return
        except ValueError:
            pass

    assert False, "Values {} and {} in path {} are different".format(
        val_out, val_ref, path
    )


def compare_json(path, out, ref, ids):
    if path in IGNORE_PATHS:
        print("Ignore path {}".format(path))
        return
    if type(ref) is str and ref == NO_TEST_VALUE:
        print("Ignore path {}".format(path))
        return
    assert type(out) is type(ref), "Incompatible type for path {}".format(path)
    if type(out) is dict:
        compare_json_dicts(path, out, ref, ids)
    elif type(out) is list:
        compare_json_lists(path, out, ref, ids)
    elif type(out) is set:
        assert False, "Unexpected type set for path {}".format(path)
    elif type(out) is tuple:
        assert False, "Unexpected type tuple for path {}".format(path)
    else:
        compare_json_values(path, out, ref, ids)


def compare_json_files(output_file_name, ref_file_name, service_config):
    json_out = None
    try:
        with open(output_file_name, "r") as out_file:
            ids = {}
            json_out = validate_ref("", None, json.load(out_file), ids, service_config)
    except FileNotFoundError:
        assert False, "Could not load file, missing output file {}".format(
            output_file_name
        )
    except Exception as e:
        assert False, "Error validating output file {}: {}".format(output_file_name, e)

    if os.getenv("OSC_GENREF", False):
        ref_exists = os.path.exists(ref_file_name)
        regenerate = True

        if ref_exists:
            with open(ref_file_name, "r") as tmp_file:
                json_ref = json.load(tmp_file)

            try:
                compare_json("", json_out, json_ref, {})
                print(
                    "Reference file {} is semantically equal, skipping regeneration".format(
                        ref_file_name
                    )
                )
                regenerate = False
            except AssertionError:
                print("Reference file {} differs, regenerating".format(ref_file_name))

        if regenerate:
            print(
                "Generating reference file {} from {}".format(
                    ref_file_name, output_file_name
                )
            )
            with open(ref_file_name, "w") as ref_file:
                ref_file.write(json.dumps(json_out, indent=4))

        return

    print("Comparing {} with {}".format(output_file_name, ref_file_name))
    try:
        with open(ref_file_name, "r") as tmp_file:
            json_ref = json.load(tmp_file)
    except FileNotFoundError:
        assert False, "Could not compare files, missing reference file {}".format(
            ref_file_name
        )
    compare_json("", json_out, json_ref, {})


def check_recreation(plan_json_path):
    try:
        with open(plan_json_path, "r") as f:
            plan = json.load(f)
    except FileNotFoundError:
        return []

    changes = []

    resource_changes = plan.get("resource_changes", [])
    for change in resource_changes:
        actions = change.get("change", {}).get("actions", [])
        resource_type = change.get("type", "unknown")
        resource_name = change.get("name", "unknown")
        path = change.get("address", f"{resource_type}.{resource_name}")

        if set(actions) == {"delete", "create"}:
            changes.append(
                {
                    "path": path,
                    "actions": actions,
                }
            )

    return changes


def create_provider_test_metaclass(root_dir, resource_filter=None):
    class ProviderTestMeta(type):
        def __new__(cls, name, bases, attrs):
            logger = logging.getLogger("tpd_test")

            def skip_tests(test_name):
                if os.path.exists("tests_to_fix.json"):
                    with open("tests_to_fix.json") as t_file:
                        skips = json.load(t_file)
                        return test_name in skips
                return False

            def create_test_func(resource, test_name, test_path):
                def func(self, tmp_path):
                    self.exec_test(test_name, test_path, tmp_path)

                func.__name__ = "test_{}_{}".format(resource, test_name)
                return func

            for resource in os.listdir(root_dir):
                if resource_filter and not resource_filter(resource):
                    continue

                path = "{}/{}".format(root_dir, resource)
                if not os.path.isdir(path):
                    logger.warning("Unexpected file: '%s'", path)
                    continue
                if resource.startswith("."):
                    continue
                for test in os.listdir(path):
                    path = "{}/{}/{}".format(root_dir, resource, test)
                    if not os.path.isdir(path):
                        logger.warning("Unexpected file: '%s'", path)
                        continue
                    logger.debug("Build test: '%s'", path)
                    func = create_test_func(resource, test, path)
                    if skip_tests(func.__name__):
                        logger.debug(
                            " %s is skipped at moment, But it must be fixed\n",
                            func.__name__,
                        )
                        continue
                    attrs[func.__name__] = func
            return type.__new__(cls, name, bases, attrs)

    return ProviderTestMeta


class BaseProviderTest:
    service_config = None

    @classmethod
    def setup_class(cls):
        cls.logger = logging.getLogger("tpd_test")
        cls.log = ""
        cls.error = False

    def setup_method(self, method):
        self.error = False
        self.work_dir = None
        self.original_dir = None
        self.log = """
==========
Log: {}
==========
        """.format(method.__name__)

    def run_cmd(self, cmd, exp_ret_code=0):
        self.logger.debug("Exec: %s", cmd)

        proc = subprocess.Popen(
            cmd,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            cwd=self.work_dir,
        )
        stdout, stderr = proc.communicate()
        stdout = stdout.decode("utf-8")
        stderr = stderr.decode("utf-8")
        if proc.returncode != exp_ret_code:
            self.error = True
            self.log += "\nERROR:\nCMD '{}' failed\nStdout: {}\nStderr: {}".format(
                cmd, stdout, stderr
            )
            print(self.log)

            assert False, "Incorrect return code {}, expected {}".format(
                proc.returncode, exp_ret_code
            )
        return stdout, stderr

    def exec_test_step(self, tf_file_path, out_file_path, is_first_step=True):
        self.logger.debug("Exec step : {}".format(tf_file_path))
        self.log += "\nTerraform validate:\n{}".format(
            self.run_cmd(["terraform validate -no-color"])[0]
        )

        plan_json_path = out_file_path.replace(".out", ".plan.json")
        plan_file_path = out_file_path.replace(".out", ".tfplan")

        self.run_cmd(
            "terraform plan -out={} -lock=false -no-color".format(plan_file_path)
        )
        self.run_cmd(
            "terraform show -json {} > {}".format(plan_file_path, plan_json_path)
        )

        if not is_first_step:
            changes = check_recreation(plan_json_path)
            if changes:
                err = "resource replacement:\n"
                for change in changes:
                    err += "  - {} (actions: {})\n".format(
                        change["path"], change["actions"]
                    )
                self.log += "\n" + err
                assert False, err

        self.log += "\nTerraform plan:\n{}".format(
            self.run_cmd("terraform plan -lock=false -no-color")[0]
        )
        self.log += "\nTerraform apply:\n{}".format(
            self.run_cmd("terraform apply -auto-approve -lock=false -no-color")[0]
        )
        self.log += "\nTerraform show:\n{}".format(
            self.run_cmd(["terraform show -no-color"])[0]
        )
        self.run_cmd(["terraform state pull > {}".format(out_file_path)])

    def exec_test(self, test_name, test_path, tmp_path):
        self.work_dir = str(tmp_path)
        self.original_dir = os.getcwd()

        test_root = os.path.dirname(__file__)
        provider_files = [
            "provider.tf",
            "variables.tf",
            "resources.auto.tfvars",
            "random.tf",
        ]
        for file_name in provider_files:
            src_file = os.path.join(test_root, file_name)
            if os.path.exists(src_file):
                shutil.copy2(src_file, self.work_dir)

        provider_dirs = ["terraform.d", "certs", "policies"]
        for dir_name in provider_dirs:
            src_dir = os.path.join(test_root, dir_name)
            if os.path.exists(src_dir):
                os.symlink(src_dir, os.path.join(self.work_dir, dir_name))

        try:
            self.logger.debug("Start test: '%s' in %s", test_name, self.work_dir)

            self.run_cmd(["terraform init -no-color"])
            stdout, _ = self.run_cmd(["terraform version -no-color"])
            self.log += "\nVERSION:{}\n".format("\n".join(stdout.splitlines()[:2]))

            tf_file_names = get_test_file_names(test_path, prefix="step", suffix=".tf")
            if not tf_file_names:
                assert False, "No step found in test directory"
            for i, tf_file_name in enumerate(tf_file_names):
                tf_file_path = os.path.join(test_path, tf_file_name)
                self.logger.debug("Process step: %s", tf_file_name)
                self.log += "\n*** step {} ***\n".format(tf_file_path)

                test_tf_path = os.path.join(self.work_dir, "test.tf")
                shutil.copy2(tf_file_path, test_tf_path)

                with open(test_tf_path) as tmp_file:
                    self.log += "\nTest file:\n{}".format(tmp_file.read())

                out_file_path = tf_file_path.replace(".tf", ".out")
                ref_file_path = tf_file_path.replace(".tf", ".ref")
                is_first_step = i == 0
                self.exec_test_step(tf_file_path, out_file_path, is_first_step)

                compare_json_files(out_file_path, ref_file_path, self.service_config)
        except Exception as error:
            self.error = True
            raise error
        finally:
            self.run_cmd("terraform destroy -auto-approve -no-color")
