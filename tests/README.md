# Outscale API provider tests

## How do test (Draft)

### Python prerequisites

You need Python 3.

You can create a Virtual Env and install dependencies following command lines below:
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt 
```

### Terraform prerequisites

You need Terraform v0.13 or greater installed in /usr/local/bin.

### Configuration

* Add AK, SK and region name in provider.auto.tfvars
* Configure account_id, image_id (and other needed resources) in resources.auto.tfvars

### Execute tests

```bash
pytest [-s] [-k <test_name>] -v ./test_provider_oapi.py
```
* use '-s' for more detailed log in console
* use '-k' to execute a subset of tests
