SHELL = /bin/bash
ROOT_PATH ?= $(PWD)
PROJECT_NAME = $$(basename $(ROOT_PATH))

include $(ROOT_PATH)/Makefile.conf

BUILD_VENV_PATH ?= "$(ROOT_PATH)/venv_build"
DEV_VENV_PATH ?= "$(ROOT_PATH)/venv_dev"

PYTHON_PATH = $$(STR="$(ROOT_PATH)"; for DEP in $(PROJECT_DEPENDENCIES); do STR=$$STR":$(ROOT_PATH)/"$$(echo $$DEP | awk -F/ '{print $$NF}'); done; echo $$STR)

help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@echo '  help          Show this help.'
	@echo '  init          Get all project dependencies.'
	@echo '  init_dev      Initialize a venv with development dependencies.'
	@echo '  init_build    Initialize a venv with dependencies needed to build package.'
	@echo '  init_ci       Add CI dependencies in development environment.'
	@echo '  pylint        Execute static code analysis.'
	@echo '  bandit        Execute static security check.'
	#@echo '  test          Execute repository tests.'
	@echo '  ci            Execute pylint, bandit and test.'
	@echo '  build         Build package.'
	@echo '  check-todo    Check remaining TODO in source code.'
	@echo '  update-deps   Update internal dependencies sha1.'
	@echo '  update-req    Update external dependencies with latest version in development environment and requirements.txt.'
	@echo ''
	@echo 'See README for more informations'
	@echo ''


python-path:
	@echo $(PYTHON_PATH)

init:
	@echo "Get dependencies from Git..."
	# install project dependencies
	for DEP in $(PROJECT_DEPENDENCIES); do                                                  \
		DEP_DIR=$(ROOT_PATH)/$$(echo $$DEP | awk -F/ '{print $$NF}');                       \
		if [ ! -d "$$DEP_DIR" ]; then                                                       \
			git clone git@gitlab.outscale.internal:qa-produit/$$DEP.git $$DEP_DIR;          \
		else                                                                                \
			pushd $$DEP_DIR; git checkout master; git pull; popd;                           \
		fi;                                                                                 \
		pushd $$DEP_DIR;                                                                    \
		git checkout $$(grep $$(echo $$DEP | awk -F/ '{print $$NF}')= $(ROOT_PATH)/internal_deps.txt | cut -d '=' -f2); \
		popd;                                                                               \
	done

init_dev: init
	@echo "Init development environment..."
	# create venv
	if [ ! -d "$(DEV_VENV_PATH)" ]; then  \
		python3 -m venv $(DEV_VENV_PATH); \
	fi
	# install python dependencies
	source $(DEV_VENV_PATH)/bin/activate; \
	pip install pip==20.0.2;              \
	pip install -U $$(ls ./*/requirements.txt | sed 's/\.\//-r \.\//g') -r requirements.txt; \
	deactivate;


init_build: init
	@echo "Init build environment..."
	# create venv
	if [ ! -d "$(BUILD_VENV_PATH)" ]; then  \
		python3 -m venv $(BUILD_VENV_PATH); \
	fi
	# install python dependencies
	source $(BUILD_VENV_PATH)/bin/activate; \
	pip install pip==20.0.2;                \
	pip install wheel==0.34.2;              \
	deactivate;

init_ci: init_dev
	@echo "Init test environment..."
	source $(DEV_VENV_PATH)/bin/activate; \
	pip install pylint==2.5.2;            \
	pip install bandit==1.6.2;            \
	pip install pytest==5.4.1;            \
	pip install pytest-cov==2.8.1;        \
	deactivate;
	#prospector==1.2.0

pylint: init_ci
	@echo "Static code analysis..."
	source $(DEV_VENV_PATH)/bin/activate;             \
	export PYTHONPATH=$(PYTHON_PATH);                 \
	pylint --rcfile=./pylint.conf $(PROJECT_MODULES); \
	[[ $$? != 0 ]] && exit -1;                        \
	deactivate;

bandit: init_ci
	@echo "Static code analysis..."
	source $(DEV_VENV_PATH)/bin/activate;          \
	export PYTHONPATH=$(PYTHON_PATH);              \
	bandit $(BANDIT_CFG) -r -l $(PROJECT_MODULES); \
	[[ $$? != 0 ]] && exit -1;                     \
	deactivate;

#test: init_ci
#	@echo "Test..."
#	source $(DEV_VENV_PATH)/bin/activate;                                                        \
#	export PYTHONPATH=$(PYTHON_PATH);                                                            \
#	export OSC_AZS=in-west-2a;                                                                   \
#	export OSC_USERS=opencell;                                                                   \
#	pytest --cov-report term --cov-report html --cov=qa_billing_tools --disable-warnings -s -v test_billing_tools; \
#	[[ $$? != 0 ]] && exit -1;                                                                   \
#	deactivate;

ci: bandit pylint# test

build: init_build
	@echo "Build package..."
	tag_name=$$(git describe --tags --exact-match 2> /dev/null);                            \
	branch_name="$$CI_BUILD_REF_NAME";                                                      \
	if [ "$$tag_name" == "" ]; then                                                         \
		tag_name="$$(git describe --tags --abbrev=4 | sed 's/-g.*//g')";                    \
		if [ "$$branch_name" == "" ]; then                                                  \
			branch_name=$$(git branch | grep \* | cut -d ' ' -f2);                          \
		fi;                                                                                 \
	else                                                                                    \
		branch_name="";                                                                     \
	fi;                                                                                     \
	echo "\"\"\"" > $(ROOT_PATH)/qa_provider_oapi/version.py;                               \
	echo "    Versionning informations" >> $(ROOT_PATH)/qa_provider_oapi/version.py;        \
	echo "    GENERATED FILE - DO NOT PUSH UPDATES" >> $(ROOT_PATH)/qa_provider_oapi/version.py; \
	echo "\"\"\"" >> $(ROOT_PATH)/qa_provider_oapi/version.py;                              \
	echo "__version__ = \"$$tag_name\"" >> $(ROOT_PATH)/qa_provider_oapi/version.py;        \
	echo "__branch__ = \"$$branch_name\"" >> $(ROOT_PATH)/qa_provider_oapi/version.py;
	source $(BUILD_VENV_PATH)/bin/activate; \
	python setup.py bdist_wheel;            \
	deactivate;

check-todo: init_ci
	@echo "Check TODO..."
	source $(DEV_VENV_PATH)/bin/activate;                     \
	export PYTHONPATH=$(PYTHON_PATH);                         \
	pylint --rcfile=./pylint.conf -e fixme $(PROJECT_MODULES) \
	deactivate;

update-deps: init
	@echo "update deps..."
	echo "# Versionning informations" > $(ROOT_PATH)/internal_deps.txt;
	for DEP in $(PROJECT_DEPENDENCIES); do                                                                                                                                     \
		DEP_DIR=$(ROOT_PATH)/$$(echo $$DEP | awk -F/ '{print $$NF}');                                                                                                          \
		pushd $$DEP_DIR;                                                                                                                                                       \
		git checkout master;                                                                                                                                                   \
		git pull;                                                                                                                                                              \
		popd;                                                                                                                                                                  \
		echo $$(echo $$DEP | awk -F/ '{print $$NF}')=$$(pushd $$DEP_DIR > /dev/null; git --no-pager log -1 --format='%H'; popd > /dev/null) >> $(ROOT_PATH)/internal_deps.txt; \
	done

update-req: init_dev
	@echo "Update requirements"
	source $(DEV_VENV_PATH)/bin/activate;                       \
	for cur in $$(grep '==' ./requirements.txt); do             \
		pkg_name=$$(echo $$cur | cut -d '=' -f1);               \
		pip install $$pkg_name --upgrade;                       \
		new=$$(pip freeze | grep $$pkg_name"=");                \
		echo "$$cur ==> $$new";                                 \
		if [ "$$(uname)" == "Darwin" ]; then                    \
			sed -i '' 's/'$$cur'/'$$new'/g' ./requirements.txt; \
		else                                                    \
			sed -i 's/'$$cur'/'$$new'/g' ./requirements.txt;    \
		fi;                                                     \
	done;                                                       \
	deactivate;
