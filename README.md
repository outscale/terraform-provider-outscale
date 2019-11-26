# Test Terraform provider for Outscale API

## Getting Started

### Python prerequisites

You will need Python 3.

You can create a Virtual Env and install dependencies following command lines below:
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt 
```

### Terraform prerequisites

You will need a Terraform v0.12.16 in your PATH.

## Usage

### Configuration

TODO: Template....

### Env setup

```bash
export OUTSCALE_OAPI_URL=outscale.com/oapi/latest
export OUTSCALE_OAPI=true
export OUTSCALE_REGION=xx-xxxxx-x
```

### Execute tests

```bash
pytest [-s] [-k <test_name>]-v ./test_provider_oapi.py
```
* use '-s' for more detailed log in console
* use '-k' to execute a subset of tests

## Add new tests

TODO...


