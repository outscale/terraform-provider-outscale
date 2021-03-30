# -*- coding:utf-8 -*-
# pylint: disable=missing-docstring

import logging
import os
import subprocess
import pytest

from check import main

ROOT_DIR = './qa_provider_oapi'

LOG_HANDLER = logging.StreamHandler()
FORMATTER = logging.Formatter('[%(asctime)s] ' +
                              '[%(levelname)8s]' +
                              '[%(module)s.%(funcName)s():%(lineno)d]: ' +
                              '%(message)s', '%m/%d/%Y %H:%M:%S')
LOG_HANDLER.setFormatter(FORMATTER)
logging.basicConfig(level=logging.DEBUG, handlers=[LOG_HANDLER])
logging.getLogger('tpd_test').setLevel(logging.DEBUG)


def get_tf_file(test_path, index):
    tf_file = None
    for tmp_file in os.listdir(test_path):
        path = "{}/{}".format(test_path, tmp_file)
        if tmp_file.startswith('step{}.'.format(index+1)):
            if tmp_file.endswith('.tf'):
                tf_file = path
    return tf_file


def get_check_files(test_path, index):
    check_files = []
    for tmp_file in os.listdir(test_path):
        path = "{}/{}".format(test_path, tmp_file)
        if tmp_file.startswith('step{}.'.format(index)):
            if tmp_file.endswith('.check'):
                check_files.append(path)
    return check_files


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
        cls.dump = True
        cls.log = None
        cls.error = False

    def setup_method(self, method):
        self.log = """
==========
Log: {}
==========
        """.format(method.__name__)
        #self.error = False
        try:
            self.run_cmd("terraform init -no-color")
            stdout, _ = self.run_cmd("terraform version -no-color")
            self.log += "\nVERSION:"
            self.log += "\n".join(stdout.splitlines()[:2])
            self.log += "\n"
        except Exception:
            try:
                self.teardown_method(method)
            except Exception:
                pass
            raise

    def teardown_method(self, method):
        try:
            pass
            #self.run_cmd("terraform destroy -force -no-color")
        finally:
            if self.error:
                self.logger.error(self.log)


    def exec_test(self, test_name, test_path):
        try:
            self.logger.debug("Start test: '%s'", test_name)
            i = 0
            if os.path.exists('{}/origin.txt'.format(test_path)):
                ret = open('{}/origin.txt'.format(test_path), 'r').read().find('WARNING')
                if ret > 0:
                    pytest.skip('WARNING during test migration')
            while True:
                try:
                    tf_file = get_tf_file(test_path, i)

                    if not tf_file:
                        break
                    i += 1
                    self.logger.debug("Process step%d: %s", i, tf_file.split('/')[-1])
                    self.log += "\n*** step {} ***\n".format(i)

                    self.run_cmd("rm -f test.tf")
                    self.run_cmd("ln -s {} test.tf".format(tf_file))

                    tmp_file = open("test.tf")
                    self.log += "\nTest file:\n"
                    self.log += tmp_file.read()
                    tmp_file.close()

                    self.logger.debug("Exec step%d", i)

                    stdout, _ = self.run_cmd("terraform validate -no-color")
                    self.log += "\nTerraform validate:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform plan -no-color")
                    self.log += "\nTerraform plan:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform apply -auto-approve -no-color")
                    self.log += "\nTerraform apply:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform show -no-color")
                    self.log += "\nTerraform show:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform state pull")
                    self.log += "\nTerraform state pull:\n"
                    self.log += stdout

                    if self.dump:
                        dump_file = tf_file.replace('.tf', '.dump')
                        tmp_file = open(dump_file, 'w')
                        tmp_file.write(stdout)
                        tmp_file.close()

                    self.logger.debug("Check step%d", i)
                    check_file_list = get_check_files(test_path, i)
                    if not check_file_list:
                        assert False, "No check file found in test directory"
                    self.run_cmd("terraform state pull > ./terraformStatePull.json")
                    for check_file in check_file_list:
                        resource = '.'.join(check_file.split('/')[-1].split('.')[1:3])
                        ret = main("./terraformStatePull.json", check_file, resource)
                        if ret:
                            self.log += "\nCheck File {}:\n".format(resource)
                            tmp_file = open(check_file)
                            self.log += tmp_file.read()
                            tmp_file.close()
                            self.log += "\n\nMissing in {}:\n".format(resource)
                            for attr in ret:
                                self.log += "  - {}\n".format(attr)
                        assert not ret


                    # TODO: add a loop for 2nd exec (ame code as 1st exec...)
                    self.logger.debug("ReExec step%d", i)

                    stdout, _ = self.run_cmd("terraform plan -no-color")
                    self.log += "\nTerraform (re)plan:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform apply -auto-approve -no-color")
                    self.log += "\nTerraform (re)apply:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform show -no-color")
                    self.log += "\nTerraform (re)show:\n"
                    self.log += stdout

                    stdout, _ = self.run_cmd("terraform state pull")
                    self.log += "\nTerraform (re)state pull:\n"
                    self.log += stdout

                    if self.dump:
                        dump_file = tf_file.replace('.tf', '.dump')
                        tmp_file = open(dump_file, 'w')
                        tmp_file.write(stdout)
                        tmp_file.close()

                    self.logger.debug("ReCheck step%d", i)
                    check_file_list = get_check_files(test_path, i)
                    if not check_file_list:
                        assert False, "No check file found in test directory"
                    self.run_cmd("terraform state pull > ./terraformStatePull.json")
                    for check_file in check_file_list:
                        resource = '.'.join(check_file.split('/')[-1].split('.')[1:3])
                        ret = main("./terraformStatePull.json", check_file, resource)
                        if ret:
                            self.log += "\nCheck File {}:\n".format(resource)
                            tmp_file = open(check_file)
                            self.log += tmp_file.read()
                            tmp_file.close()
                            self.log += "\n\nMissing in {}:\n".format(resource)
                            for attr in ret:
                                self.log += "  - {}\n".format(attr)
                        assert not ret

                finally:
                    pass
                    #self.run_cmd("rm -f test.tf")
            if i == 0:
                assert False, "No step found in test directory"
        except Exception as error:
            self.error = True
            raise
        finally:
            try:
                self.run_cmd("terraform destroy -force -no-color")
            finally:
                self.run_cmd("rm -f test.tf")
                self.run_cmd("rm -f terraform.tfstate")


    def run_cmd(self, cmd, exp_ret_code=0):
        self.logger.debug("Exec: %s", cmd)
        proc = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = proc.communicate()
        stdout = stdout.decode("utf-8")
        stderr = stderr.decode("utf-8")
        #if self.debug or proc.returncode != exp_ret_code:
        #    self.logger.debug("stdout:\n%s", stdout)
        #    self.logger.debug("stderr:\n%s", stderr)
        if proc.returncode != exp_ret_code:
            self.error = True
            self.log += "\nERROR:"
            self.log += "\nCMD '{}' failed".format(cmd)
            self.log += "\nStdout: "
            self.log += stdout
            self.log += "\nStderr: "
            self.log += stderr
            assert proc.returncode == exp_ret_code
        return (stdout, stderr)
