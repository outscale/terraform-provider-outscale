# -*- coding:utf-8 -*-
# pylint: disable=missing-docstring

import logging
import os
import subprocess
import pytest
import json

from qa_provider_oapi.check import main

# TODO check if this could be better using regular expressions
IGNORE_PATHS = ['.lineage',
                '.serial',
                '.terraform_version',
                '.version',
                '.resources.provider']
TAG_END_PATHS = ['.id', '_id']
NO_TEST_VALUE = '########'
IGNORE_END_ELEMENTS = ['request_id',
                       'mac_address',
                       'public_ip',
                       'keypair_fingerprint',
                       'link_date',
                       'private_key',
                       'private_ip',
                       'filter',
                       'image_name',
                       'creation_date',
                       'public_dns_name',
                       'private_dns_name',
                       'dns_name',
    		       'secret_key',
                       'cookie_name',
                       'client_gateway_configuration',
                       'last_modification_date',
                       'upload_date' ,
                       'comment',
                       'osu_manifest_url',
                       'max_value',
                       'file_location',
                       'message',
                       'ip_ranges']
IGNORE_END_PATHS = []
TINA_ID_PREFIXES = ['i', 'subnet', 'snap', 'img', 'vol', 'eni', 'vpc', 'igw', 'nat', 'vgw', 'pcx', 'sg', 'rtb', 'rtbassoc', 'vpn', 'vpcconn', 'ami', 'dxvif','vpce','fgpu']
VARIABLES_FILE_NAME = ['provider.auto.tfvars', 'resources.auto.tfvars']
VARIABLES = ['region']
SET_KEY_VALUES = ['resources', 'tags']
ID_PREFIX = '##id-'
ID_SUFFIX = '##'

ROOT_DIR = os.path.join(os.path.dirname(__file__), 'data')
LOG_HANDLER = logging.StreamHandler()
FORMATTER = logging.Formatter('[%(asctime)s] ' +
                              '[%(levelname)8s]' +
                              '[%(module)s.%(funcName)s():%(lineno)d]: ' +
                              '%(message)s', '%m/%d/%Y %H:%M:%S')
LOG_HANDLER.setFormatter(FORMATTER)
logging.basicConfig(level=logging.DEBUG, handlers=[LOG_HANDLER])
logging.getLogger('tpd_test').setLevel(logging.DEBUG)

terraform_vars = {}
for file_name in VARIABLES_FILE_NAME:
    file_name = os.path.abspath(os.path.join(os.path.dirname(__file__), file_name))
#    print(file_name)
    with open(file_name, 'r') as var_file:
        lines = var_file.readlines()
        for line in lines:
            line = line.strip()
            if line.startswith('#'):
                continue
            elts = line.split('=')
            if len(elts) != 2:
                continue
            terraform_vars[elts[0].strip()] = elts[1].strip().strip("\"")


def get_test_file_names(test_path, prefix='step', suffix='.tf'):
    ret_file_names = []
    for tmp_file in os.listdir(test_path):
        if tmp_file.startswith(prefix) and tmp_file.endswith(suffix):
            ret_file_names.append(tmp_file)
    return sorted(ret_file_names)


def validate_value_ref(path, parent, value, ids):

    replace_value = None
    replace = None
    if type(value) == str:
        value_items = value.split('-')
        if len(value_items) == 2 and len(value_items[1]) == 8 and value_items[0] in TINA_ID_PREFIXES:
            replace = value

    for p in TAG_END_PATHS:
        if path.endswith(p):
            replace = value
            break
    if replace:
        if replace not in ids:
            ids[replace] = '{}{}{}'.format(ID_PREFIX, len(ids.keys()), ID_SUFFIX)
        replace_value = ids[replace]

    if not replace_value and type(value) == str:
        tmp_value = value
        for var in terraform_vars:
            if terraform_vars[var] in tmp_value:
                tmp_value = tmp_value.replace(terraform_vars[var], '###{}###'.format(var))
        if tmp_value != value:
            replace_value = tmp_value

    if replace_value:
        path_end = path.split('.')[-1]
        if type(parent) == list:
            parent[int(path_end)] = replace_value
        else:
            parent[path_end] = replace_value


def validate_dict_ref(path, json_ref, ids):
    for key, value in json_ref.items():
        validate_ref('{}.{}'.format(path, key), json_ref, value, ids)


def validate_list_ref(path, json_ref, ids):
    for i in range(len(json_ref)):
        validate_ref('{}.{}'.format(path, i), json_ref, json_ref[i], ids)


def validate_ref(path, parent, value, ids):
    path_end = path.split('.')[-1]
    if path in IGNORE_PATHS or path_end in IGNORE_END_ELEMENTS:
        parent[path_end] = NO_TEST_VALUE
        return
    for p in IGNORE_END_PATHS:
        if path.endswith(p):
            parent[path_end] = NO_TEST_VALUE
            return

    if type(value) == dict:
        validate_dict_ref(path, value, ids)
    elif type(value) == list:
        validate_list_ref(path, value, ids)
    elif type(value) == tuple:
        assert False, 'Unexpected type tuple for path {}'.format(path)
    else:
        validate_value_ref(path, parent, value, ids)
    return value


def compare_json_dicts(path, dict_out, dict_ref, ids):
    keys_out = sorted(set(dict_out.keys()))
    keys_ref = sorted(set(dict_ref.keys()))
    assert len(keys_out) == len(keys_ref), 'Not the same keys number for path {}'.format(path)
    for key in keys_out:
        assert key in keys_ref, 'Could not find key {}.{} in output'.format(path, key)
    for key in keys_ref:
        assert key in keys_out, 'Could not find key {}.{} in reference'.format(path, key)
    for key in keys_out:
        do_set = False
        if key in SET_KEY_VALUES:
            do_set = True
        elif key.endswith('s') and type(dict_out[key]) == list and dict_out[key] and type(dict_out[key][0]) == dict and '{}_id'.format(key[:-1]) in dict_out[key][0]:
            do_set = True
        if do_set:
            compare_json_sets('{}.{}'.format(path, key), dict_out[key], dict_ref[key], ids)
        else:
            compare_json('{}.{}'.format(path, key), dict_out[key], dict_ref[key], ids)


#def compare_json_lists(path, list_out, list_ref, ids):
#    assert len(list_out) == len(list_ref)
#    for i in range(len(list_out)):
#        compare_json('{}.{}'.format(path, i), list_out[i], list_ref[i], ids)

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
                compare_json('{}'.format(path), out_elt, ref_elt, tmp_ids)
                found_elts.append(ref_elt)
                current_ids = tmp_ids
                errors = []
                break
            except Exception as error:
                tmp_ids = current_ids
                errors.append(error)
                pass
        if errors:
            assert False, 'Could not match set values for path {}, {}'.format(path, errors)
    ids = current_ids 

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
                compare_json('{}'.format(path), out_elt, ref_elt, tmp_ids)
                found_elts.append(ref_elt)
                current_ids = tmp_ids
                errors = []
                break
            except Exception as error:
                tmp_ids = current_ids
                errors.append(error)
                pass
        if errors:
            assert False, 'Could not match set values for path {}, {}'.format(path, errors)
    ids = current_ids


def compare_json_values(path, val_out, val_ref, ids):
    # do not check values that should be ignored
    if val_ref == NO_TEST_VALUE:
        return

    if val_out == val_ref:
        return

    # accept id change
    if val_out in ids and ids[val_out] == val_ref:
        return

    if val_out not in ids and val_out.startswith(ID_PREFIX) and val_ref.startswith(ID_PREFIX):
        ids[val_out] = val_ref
        return

    # accept string values that only differ by the last digit
    if type(val_out) == str and len(val_out) > 1 and len(val_ref) > 1 and val_out[:-1] == val_ref[:-1]:
        try:
            int(val_out[-1])
            int(val_ref[-1])
            return
        except ValueError:
            pass

    try:
        assert False, 'Values {} and {} in path {} are different'.format(val_out, val_ref, path)
    except AssertionError as error:
        raise error


def compare_json(path, out, ref, ids):
    if path in IGNORE_PATHS:
        print('Ignore path {}'.format(path))
        return
    if type(ref) == str and ref == NO_TEST_VALUE:
        print('Ignore path {}'.format(path))
        return
    assert type(out) == type(ref), 'Incompatible type for path {}'.format(path)
    if type(out) == dict:
        compare_json_dicts(path, out, ref, ids)
    elif type(out) == list:
        compare_json_lists(path, out, ref, ids)
    elif type(out) == set:
        assert False, 'Unexpected type set for path {}'.format(path)
    elif type(out) == tuple:
        assert False, 'Unexpected type tuple for path {}'.format(path)
    else:
        compare_json_values(path, out, ref, ids)


def compare_json_files(output_file_name, ref_file_name):
    json_out = None
    try:
        with open(output_file_name, 'r') as out_file:
            ids = {}
            json_out = validate_ref('', None, json.load(out_file), ids)
    except FileNotFoundError:
        assert False, 'Could load file, missing output file {}'.format(ref_file_name)

    if os.getenv('OSC_GENREF', False):
        print('Generating reference file {} from {}'.format(ref_file_name, output_file_name))
        with open(ref_file_name, 'w') as ref_file:
            ref_file.write(json.dumps(json_out, indent=4))
        return

    print('Comparing {} with {}'.format(output_file_name, ref_file_name))
    try:
        with open(ref_file_name, 'r') as tmp_file:
            json_ref = json.load(tmp_file)
    except FileNotFoundError:
        assert False, 'Could not compare files, missing reference file {}'.format(ref_file_name)
    compare_json('', json_out, json_ref, {})


class ProviderOapiMeta(type):
    def __new__(cls, name, bases, attrs):
        logger = logging.getLogger('tpd_test')

        def create_test_func(resource, test_name, test_path):
            def func(self):
                self.exec_test(test_name, test_path)

            func.__name__ = "test_{}_{}".format(resource, test_name)
            return func

        for resource in os.listdir(ROOT_DIR):
            path = "{}/{}".format(ROOT_DIR, resource)
            if not os.path.isdir(path):
                logger.warning("Unexpected file: '%s'", path)
                continue
            if resource.startswith('.'):
                continue
            for test in os.listdir(path):
                path = "{}/{}/{}".format(ROOT_DIR, resource, test)
                if not os.path.isdir(path):
                    logger.warning("Unexpected file: '%s'", path)
                    continue
                logger.debug("Build test: '%s'", path)
                func = create_test_func(resource, test, path)
                attrs[func.__name__] = func
        return type.__new__(cls, name, bases, attrs)


class TestProviderOapi(metaclass=ProviderOapiMeta):

    @classmethod
    def setup_class(cls):
        cls.logger = logging.getLogger('tpd_test')
        cls.log = None
        cls.error = False

    def setup_method(self, method):
        self.log = """
==========
Log: {}
==========
        """.format(method.__name__)
        self.error = False
        try:
            self.run_cmd("terraform init -no-color")
            stdout, _ = self.run_cmd("terraform version -no-color")
            self.log += "\nVERSION:{}\n".format("\n".join(stdout.splitlines()[:2]))
        except Exception:
            try:
                self.teardown_method(method)
            except Exception:
                pass
            raise

    def teardown_method(self, method):
        try:
            pass
            # self.run_cmd("terraform destroy -force -no-color")
        finally:
            if self.error:
                self.logger.error(self.log)

    def run_cmd(self, cmd, exp_ret_code=0):
        self.logger.debug("Exec: %s", cmd)
        proc = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = proc.communicate()
        stdout = stdout.decode("utf-8")
        stderr = stderr.decode("utf-8")
        # if self.debug or proc.returncode != exp_ret_code:
        # self.logger.debug("stdout:\n%s", stdout)
        # self.logger.debug("stderr:\n%s", stderr)
        if proc.returncode != exp_ret_code:
            self.error = True
            self.log += "\nERROR:\nCMD '{}' failed\nStdout: {}\nStderr: {}".format(cmd, stdout, stderr)
            print(self.log)

            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            
            assert False, 'Incorrect return code {}, expected {}'.format(proc.returncode, exp_ret_code)
        return stdout, stderr

    def exec_test_step(self, tf_file_path, out_file_path):
        self.logger.debug("Exec step : {}".format(tf_file_path))
        self.log += "\nTerraform validate:\n{}".format(self.run_cmd("terraform validate -no-color")[0])
        self.log += "\nTerraform plan:\n{}".format(self.run_cmd("terraform plan -lock=false -no-color")[0])
        self.log += "\nTerraform apply:\n{}".format(self.run_cmd("terraform apply -auto-approve -lock=false -no-color")[0])
        self.log += "\nTerraform show:\n{}".format(self.run_cmd("terraform show -no-color")[0])
        self.run_cmd("terraform state pull > {}".format(out_file_path))

    def exec_test(self, test_name, test_path):
        try:
            self.logger.debug("Start test: '%s'", test_name)
            if os.path.exists('{}/origin.txt'.format(test_path)):
                ret = open('{}/origin.txt'.format(test_path), 'r').read().find('WARNING')
                if ret > 0:
                    pytest.skip('WARNING during test migration')
            tf_file_names = get_test_file_names(test_path, prefix='step', suffix='.tf')
            # TDOD erase once not needed anymore
            check_file_names = get_test_file_names(test_path, prefix='step', suffix='.check')
            check_file_index = 0
            if not tf_file_names:
                assert False, "No step found in test directory"
            for tf_file_name in tf_file_names:
                try:
                    tf_file_path = os.path.join(test_path, tf_file_name)
                    self.logger.debug("Process step: %s", tf_file_name)
                    self.log += "\n*** step {} ***\n".format(tf_file_path)

                    self.run_cmd("rm -f test.tf")
                    self.run_cmd("ln -s {} test.tf".format(tf_file_path))

                    with open("test.tf") as tmp_file:
                        self.log += "\nTest file:\n{}".format(tmp_file.read())

                    out_file_path = tf_file_path.replace('.tf', '.out')
                    ref_file_path = tf_file_path.replace('.tf', '.ref')
                    self.exec_test_step(tf_file_path, out_file_path)

                    # TDOD erase once not needed anymore
                    if os.getenv('OSC_GENREF', False) or os.path.exists(ref_file_path):
                        compare_json_files(out_file_path, ref_file_path)
                    else:
                        check_file_name = check_file_names[check_file_index]
                        resource = '.'.join(check_file_name.split('.')[1:3])
                        ret = main(out_file_path, os.path.join(test_path, check_file_name), resource)
                        if ret:
                            self.log += "\nCheck File {}:\n".format(resource)
                            with open(os.path.join(test_path, check_file_name)) as tmp_file:
                                self.log += tmp_file.read()
                            self.log += "\n\nMissing in {}:\n".format(resource)
                            for attr in ret:
                                self.log += "  - {}\n".format(attr)
                        assert not ret
                finally:
                    check_file_index += 1
                    pass
        except Exception as error :
            self.error = True
            raise error
        finally:
            try:
                self.run_cmd("terraform destroy -auto-approve -no-color")
            finally:
                self.run_cmd("rm -f test.tf")
                self.run_cmd("rm -f terraform.tfstate")
                self.run_cmd("rm -f .terraform.lock.hcl")
