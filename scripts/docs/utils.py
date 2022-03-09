import re

def camel_case_to_snake_case(value):
    """Convert from camel case to snake case."""
    # re.sub works on non overlapping occurencies only.
    s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', value)
    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()

def snake_case_to_camel_case(value):
    return ''.join([x.title() for x in value.split('_') or []])