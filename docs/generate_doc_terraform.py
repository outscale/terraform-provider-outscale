#!/usr/bin/env python3
# -*- coding:utf-8 -*-

import argparse
import csv
import io
import re
import shutil
import os
import yaml

import gtw.api.openapi.parser as openapi_parser
import gtw.api.schema as schema
import gtw.api.utils as utils

# NOTE : python3 generate_doc_terraform.py --provider_directory /Users/mercatguillaume/workspace/terraform-provider-outscale/outscale --api /Users/mercatguillaume/workspace/osc-api/outscale.yaml --output_directory /Users/mercatguillaume/workspace/output_test --template_directory /Users/mercatguillaume/workspace/api-framework/terraform_template

parser = argparse.ArgumentParser(description='Generate documentation terraform')
parser.add_argument('--new_format', action='store_true')
parser.add_argument('--no-addapt', action='store_true')
parser.add_argument('--provider_directory', required=True)
parser.add_argument('--template_directory', required=True)
parser.add_argument('--api', help='Source oAPI specification path.',
                    required=True)
parser.add_argument('--output_directory', help='Output directory path.',
                    required=True)
ARGS = parser.parse_args()


def print_dict(item, path, profondeur):
    result = str()
    if isinstance(item, schema.ObjectField):
        # if path not in code_file:
        #    print('    - {} not found'.format(path))
        result += '    ' * profondeur + '* `{}` - {}\n'.format(path,
                                                            item.description)
        profondeur += 1
        for x, y in sorted(item.properties.items()):
            next_path = utils.camel_case_to_snake_case(x)
            if isinstance(y, schema.TerminalField):
                # if next_path not in code_file:
                #    print('    - {} not found'.format(next_path))
                result += '    ' * profondeur + '* `{}` - {}\n'.format(next_path, y.description)
            else:
                result += print_dict(
                    # code_file,
                    y, next_path, profondeur)
        return result
    elif isinstance(item, schema.ArrayField):
        if isinstance(item.item, schema.ObjectField):
            #if path not in code_file:
            #    print('    - {} not found'.format(path))
            result += '    ' * profondeur + '* `{}` - {}\n'.format(path, item.description)
            profondeur += 1
            for x, y in sorted(item.item.properties.items()):
                next_path = utils.camel_case_to_snake_case(x)
                if isinstance(y, schema.TerminalField):
                    #if next_path not in code_file:
                    #    print('{} not found'.format(next_path))
                    result += '    ' * profondeur + '* `{}` - {}\n'.format(next_path, y.description)
                else:
                    result += '{}'.format(print_dict(
                        #code_file,
                        y, next_path, profondeur))
            return result
        elif isinstance(item.item, schema.TerminalField):
            #if path not in code_file:
            #    print('    - {} not found'.format(path))
            return '    ' * profondeur + '* `{}` - {}\n'.format(path, item.description)
    elif isinstance(item, schema.TerminalField):
        #if path not in code_file:
        #    print('    - {} not found'.format(path))
        return '    ' * profondeur + '* `{}` - {}\n'.format(path, item.description)
    else:
        return path + '\n'


def print_dict_input(item, path, profondeur, mandatories=[]):
    result = str()
    if isinstance(item, schema.ObjectField):
        #if path not in code_file:
        #    print('    - {} not found'.format(path))
        result += '    ' * profondeur + '* `{}` - {}\n'.format(path,
                                                              item.description)
        profondeur += 1
        for x, y in sorted(item.properties.items()):
            next_path = utils.camel_case_to_snake_case(x)
            if isinstance(y, schema.TerminalField):
                #if next_path not in code_file:
                #    print('    - {} not found'.format(next_path))
                result += '    ' * profondeur + '* `{}` - ({}) {}\n'.format(next_path,
                                                      'Required'
                                                      if y.name in mandatories
                                                      else 'Optional',
                                                      y.description)
            else:
                result_calculed = print_dict_input(
                    #code_file,
                    y,
                    next_path,
                    profondeur,
                    mandatories)
                result += '{}'.format(result_calculed)
        return result
    elif isinstance(item, schema.ArrayField):
        if isinstance(item.item, schema.ObjectField):
            #if path not in code_file:
            #    print('    - {} not found'.format(path))
            result += '    ' * profondeur + '* `{}` - ({}) {}\n'.format(path,
                                                  'Required'
                                                  if item.name in mandatories
                                                  else 'Optional',
                                                  item.description)
            profondeur += 1
            for x, y in sorted(item.item.properties.items()):
                next_path = utils.camel_case_to_snake_case(x)
                if isinstance(y, schema.TerminalField):
                    #if next_path not in code_file:
                    #    print('    - {} not found'.format(next_path))
                    result += '    ' * profondeur + '* `{}` - ({}) {}\n'.format(next_path,
                                                          'Required'
                                                          if y.name in mandatories
                                                          else 'Optional',
                                                          y.description)
                else:
                    result_calculed = print_dict_input(
                        #code_file,
                        y, next_path,
                        profondeur, mandatories)
                    result += '{}'.format(result_calculed)
            return result
        elif isinstance(item.item, schema.TerminalField):
            #if path not in code_file:
            #    print('    - {} not found'.format(path))
            return '    ' * profondeur + '* `{}` - ({}) {}\n'.format(path,
                                               'Required'
                                               if item.name in mandatories
                                               else 'Optional',
                                               item.description)
    elif isinstance(item, schema.TerminalField):
        #if path not in code_file:
        #    print('    - {} not found'.format(path))
        return '    ' * profondeur + '* `{}` - ({}) {}\n'.format(path,
                                           'Required'
                                           if item.name in mandatories
                                           else 'Optional',
                                           item.description)
    else:
        return path + '\n'


def addField(field_object, key, value):
    print('add k: {}, v: {}'.format(key, value))
    if '.' in key:
        key_split = key.split('.', 1)
        if isinstance(field_object, schema.ObjectField):
            if key_split[0] in field_object.properties:
                addField(field_object.properties[key_split[0]], key_split[1], value)
        else:
            if key_split[0] in field_object:
                addField(field_object[key_split[0]], key_split[1], value)
    else:
        new_obj = schema.ObjectField({}, key, 'object', 'in', 'toto', description=value)
        if isinstance(field_object, schema.ObjectField):
            print('ObjectField')
            field_object.properties[key] = new_obj
        elif isinstance(field_object, schema.ArrayField):
            field_object.item.properties[key] = new_obj
        else:
            print('dict')
            field_object[key] = new_obj


def renameField(field_object, key, value):
    if '.' in key:
        key_split = key.split('.', 1)
        renameField(field_object[key_split[0]], key_split[1], value)
    else:
        if isinstance(field_object, schema.ObjectField):
            if key in field_object.properties.keys():
                old = field_object.properties[key]
                if isinstance(value, str):
                    new_key = utils.snake_case_to_camel_case(value)
                    field_object.properties[new_key] = old
                    del field_object.properties[key]
                else:
                    new_key = value.get('name', key)
                    old.description = value.get('description', old.description)
                    if utils.snake_case_to_camel_case(new_key) != key:
                        field_object.properties[utils.snake_case_to_camel_case(new_key)] = old
                        del field_object.properties[key]
            else:
                print('[WARNING]: you want to rename an unknown field "{}"'.format(key))
        elif isinstance(field_object, schema.ArrayField):
            if key in field_object.item.properties.keys():
                old = field_object.item.properties[key]

                if isinstance(value, str):
                    new_value = utils.snake_case_to_camel_case(value)
                    field_object.item.properties[new_value] = old
                    del field_object.item.properties[key]
                else:
                    new_key = value.get('name', key)
                    old.description = value.get('description', old.description)
                    if utils.snake_case_to_camel_case(new_key) != key:
                        field_object.item.properties[utils.snake_case_to_camel_case(new_key)] = old
                        del field_object.item.properties[key]
            else:
                print('[WARNING]: you want to rename an unknown field "{}"'.format(key))
        else:
            if key in field_object:
                old = field_object[key]

                if isinstance(value, str):
                    new_value = utils.snake_case_to_camel_case(value)
                    field_object[new_value] = old
                    del field_object[key]
                else:
                    new_key = value.get('name', key)
                    old.description = value.get('description', old.description)
                    if utils.snake_case_to_camel_case(new_key) != key:
                        field_object[utils.snake_case_to_camel_case(new_key)] = old
                        del field_object[key]
            else:
                print('[WARNING]: you want to rename an unknown field "{}"'.format(key))


def removeField(field_object, key):
    if '.' in key:
        key_split = key.split('.', 1)
        removeField(field_object[key_split[0]], key_split[1])
    else:
        if isinstance(field_object, schema.ObjectField):
            if key in field_object.properties.keys():
                del field_object.properties[key]
            else:
                print('[WARNING]: you want to delete an unknown field "{}"'.format(key))
        elif isinstance(field_object, schema.ArrayField):
            if key in field_object.item.properties.keys():
                del field_object.item.properties[key]
            else:
                print('[WARNING]: you want to delete an unknown field "{}"'.format(key))
        else:
            if key in field_object:
                del field_object[key]
            else:
                print('[WARNING]: you want to delete an unknown field "{}"'.format(key))


def treatAddPropData(field_to_update, part_to_update, addprop_content):
    arg = addprop_content.get(part_to_update, None)    
    if arg and field_to_update:
        for k, v in arg.get('add', {}).items():
            new_k = utils.snake_case_to_camel_case(k)
            addField(field_to_update, new_k, v)
        for k, v in arg.get('rename', {}).items():
            new_k = utils.snake_case_to_camel_case(k)
            renameField(field_to_update, new_k, v)
        for k in arg.get('remove', []):
            new_k = utils.snake_case_to_camel_case(k)
            removeField(field_to_update, new_k)


def file_template(template, links, resource_name, is_singular, input_field, output_field, example_content, import_content):
    resource_name_singular = resource_name if is_singular else resource_name[:-1]
    placeholders = links.get(resource_name_singular, 'NOT_FOUND {}'.format(resource_name_singular))
    print('======= placeholders -> {}'.format(placeholders))

    content_file = template.replace(
        'LINK_UG',
        links.get(resource_name_singular, {}).get('LINK_UG', 'NOT_FOUND'))
    content_file = content_file.replace(
        'LINK_API',
        links.get(resource_name_singular, {}).get('LINK_API', 'NOT_FOUND'))

    content_file = content_file.replace(
        'LITERAL_NAME_PLURAL',
        links.get(resource_name_singular, {}).get('LITERAL_NAME_PLURAL', 'NOT_FOUND'))

    if 's' == resource_name[-1:]:
        content_file = content_file.replace(
            'ARTICLE ', '')
        content_file = content_file.replace(
            'LITERAL_NAME',
            links.get(resource_name_singular, {}).get('LITERAL_NAME_PLURAL', 'NOT_FOUND'))
    else:
        content_file = content_file.replace(
            'ARTICLE', links.get(resource_name_singular, {}).get('ARTICLE', 'NOT_FOUND'))
        content_file = content_file.replace(
            'LITERAL_NAME',
            links.get(resource_name_singular, {}).get('LITERAL_NAME', 'NOT_FOUND'))

    content_file = content_file.replace('RESOURCE_NAME', resource_name)
    content_file = content_file.replace('RESOURCE-NAME',
                                        resource_name.replace('_', '-'))

    content_file = content_file.replace('ARGUMENTS_SENTENCE', 'The following arguments are supported:' if input_field else 'No argument is supported.')
    content_file = content_file.replace('INPUT' if input_field else 'INPUT\n', input_field)

    content_file = content_file.replace('ATTRIBUTES_SENTENCE', 'The following attributes are exported:' if output_field else  'No attribute is exported.')
    content_file = content_file.replace('OUTPUT' if output_field else 'OUTPUT\n', output_field)

    full_example_content = example_content if len(example_content) else ''
    content_file = content_file.replace('EXAMPLE' if full_example_content else 'EXAMPLE\n', full_example_content)
    content_file = content_file.replace('IMPORT', import_content)

    content_file = content_file.replace('](#', '](https://docs.outscale.com/api#')

    return content_file


def main():
    provider_filename = '{}/provider.go'.format(ARGS.provider_directory)
    with io.open(provider_filename, 'r') as f:
        provider_file = f.read()

    with io.open('{}/template_ressource.md'.format(ARGS.template_directory),
                 'r') as f:
        template_resource = f.read()

    with io.open('{}/template_datasource.md'.format(ARGS.template_directory),
                 'r') as f:
        template_datasource = f.read()

    with io.open('{}/template_datasources.md'.format(ARGS.template_directory),
                 'r') as f:
        template_datasources = f.read()

    resources = {}
    with io.open('{}/resources.csv'.format(ARGS.template_directory), 'r',
                 newline='', encoding='utf-8') as csv_file:
        values = csv.reader(csv_file, delimiter=',')
        if not values:
            print('No data found.')
        else:
            for row in values:
                resources[row[0]] = [x[1:] for x in row[1:]]

    links = {}
    with io.open('{}/links.csv'.format(ARGS.template_directory), 'r',
                 newline='', encoding='utf-8') as csv_file:
        values = csv.reader(csv_file, delimiter=',')
    
        if not values:
            print('No data found.')
        else:
            for row in values:
                links[row[0]] = {
                    'ARTICLE': row[1],
                    'LITERAL_NAME': row[2],
                    'LITERAL_NAME_PLURAL': row[3],
                    'LINK_UG': row[4],
                    'LINK_API': row[5],
                }

    extention = '.md' if ARGS.new_format else '.html.markdown'
    index_dirpath = ARGS.output_directory + ('/docs'
                                               if ARGS.new_format else '/website/docs')
    navbar_dirpath = ARGS.output_directory + ('/docs'
                                               if ARGS.new_format else '/website')
    navbar_file = """<% wrap_layout :inner do %>
  <% content_for :sidebar do %>
    <h4>OUTSCALE</h4>

    <ul class="nav docs-sidenav">
      <li>
        <a href="#">Data Sources</a>
        <ul class="nav">
"""
    navbar_data_source={}
    navbar_resource={}
    print('Parsing API from {}...'.format(ARGS.api))
    oapi = openapi_parser.parse(ARGS.api)

    for name, call_list in resources.items():
        filename = name
        input_fields = set()
        output_fields = set()
        dirpath = str()
        template = str()
        example_content = str()
        import_content = str()
        addprop_content = {}
        resource_name = str()
        code_filename = '{}/{}'.format(
            ARGS.provider_directory, filename[len('outscale/')+1:])
        is_singular = True

        if 'data_source' in name:
            dirpath = ARGS.output_directory + (
                '/docs/data-sources' if ARGS.new_format else '/website/docs/d')
            resource_name = re.search('outscale/data_source_outscale_(.*).go', filename).group(1)
            navbar_data_source[resource_name] = '/docs/providers/outscale/d/{}.html'.format(resource_name)
            template = template_datasource
            # Load example, import and addprop
            try:
                with io.open('{}/Content/datasources/{}-example.md'.format(ARGS.template_directory, resource_name),
                             'r') as f:
                    example_content = f.read()
            except FileNotFoundError as e:
                pass
            try:
                with io.open('{}/Content/datasources/{}-import.md'.format(ARGS.template_directory, resource_name),
                             'r') as f:
                    import_content = f.read()
            except FileNotFoundError as e:
                pass
            try:
                with io.open('{}/Content/datasources/{}-addprop.yaml'.format(ARGS.template_directory, resource_name),
                             'r') as f:
                    addprop_content = yaml.load(f, yaml.FullLoader)
            except FileNotFoundError as e:
                pass
            if resource_name == 'vms_state':
                template = template_datasources
                is_singular = False
                resource_name = 'vm_states'
            elif resource_name[-1] == 's':
                template = template_datasources
                is_singular = False
        elif 'resource' in name:
            dirpath = ARGS.output_directory + (
                '/docs/resources' if ARGS.new_format else '/website/docs/r')
            resource_name = re.search('outscale/resource_outscale_(.*).go', filename).group(1)
            navbar_resource[resource_name] = '/docs/providers/outscale/r/{}.html'.format(resource_name)
            template = template_resource
            # Load example, import and addprop
            try:
                with io.open('{}/Content/resources/{}-example.md'.format(ARGS.template_directory, resource_name),
                             'r') as f:
                    example_content = f.read()
            except FileNotFoundError as e:
                pass
            try:
                with io.open('{}/Content/resources/{}-import.md'.format(ARGS.template_directory, resource_name),
                             'r') as f:
                    import_content = f.read()
            except FileNotFoundError as e:
                pass
            try:
                with io.open('{}/Content/resources/{}-addprop.yaml'.format(ARGS.template_directory, resource_name),
                             'r') as f:
                    addprop_content = yaml.load(f, yaml.FullLoader)
            except FileNotFoundError as e:
                pass
        else:
            print('This filename, {} is not in a known format - we do not treat it.'.format(name))
            continue

        if 'outscale_{}'.format(resource_name) not in provider_file:
            print('{} not found'.format('outscale_{}'.format(resource_name)))
            continue

        with io.open(code_filename, 'r') as f:
            code_file = f.read()

        print('\nTreating {} '.format(filename))
        call_complete = schema.Call('', {}, {}, [], '', False)
        for call in call_list:
            call_complete.merge(oapi.calls.get(call))
            oapi_call = oapi.calls.get(call)
            if not oapi_call:
                print('This call is not found : {}'.format(call))
                continue

        print('Outscale api call found: {}'.format(call))

        print(' - Remove double (Sing/Plur) in output fileds')
        for k in list(call_complete.output_fields.keys()):
            if k[-1] == 's' and k[:-1] in call_complete.output_fields.keys():
                if is_singular:
                    print('del {}'.format(k))
                    del call_complete.output_fields[k]
                else:
                    print('del {}'.format(k[:-1]))
                    del call_complete.output_fields[k[:-1]]
        
        print(' - Change output fields format for singular data source and resources')
        if is_singular or 'resource' in name:
            update_fields = dict()
            for k in list(call_complete.output_fields.keys()):
                if k == 'ResponseContext':
                    update_fields[k] = call_complete.output_fields[k]
                    continue
                print("Exploring {}".format(k))
                value = call_complete.output_fields[k]
                print("What ?  {}, Type ? {}".format(k, type(value)))
                if isinstance(value, schema.ArrayField):
                    if isinstance(value.item, schema.ObjectField):
                        for inner_k in value.item.properties.keys():
                            update_fields[inner_k] = value.item.properties[inner_k]
                    elif isinstance(value.item, schema.TerminalField):
                        update_fields[inner_k] = value.item[inner_k]
                    else:
                        print("Who are you ? ... What: {} Type:{}".format(inner_k, type(value.item)))
                        exit(1)
                elif isinstance(value, schema.ObjectField):
                    for inner_k in value.properties.keys():
                        update_fields[inner_k] = value.properties[inner_k]
                elif isinstance(value, schema.TerminalField):
                    update_fields[k] = value
                else:
                    print("Strange ... What : {} Type:{}".format( k, type(value)))
                    exit(1)
            call_complete.output_fields = update_fields
        
        print(' - Treating addprop data')
        input_field_to_update = call_complete.input_fields if call_complete else {}
        if not ARGS.no_addapt:
            print("Exploring argument addprop")
            treatAddPropData(input_field_to_update, 'argument', addprop_content)
        output_field_to_update = call_complete.output_fields if call_complete else {}
        if not ARGS.no_addapt:
            print("Exploring attributes addprop")
            treatAddPropData(output_field_to_update, 'attribute', addprop_content)
        
        print(' - Treating input parameters ...')
        input_field_to_parse = call_complete.input_fields if call_complete else {}
        for a, b in input_field_to_parse.items():
            if a != 'DryRun':
                input_fields.add(print_dict_input(
                        b,
                        utils.camel_case_to_snake_case(a),
                        0,
                        call_complete.required))

        print(' - Treating output parameters ...')
        output_field_to_parse = call_complete.output_fields if call_complete else {}
        for a, b in output_field_to_parse.items():
            if a != 'ResponseContext':
                output_fields.add(print_dict(
                        b,
                        utils.camel_case_to_snake_case(a),
                        0))

        str_input = str()
        for x in sorted(input_fields):
            str_input += x
        str_output = str()
        for x in sorted(output_fields):
            str_output += x
        content_file = file_template(template,
                                     links,
                                     resource_name,
                                     is_singular,
                                     str_input,
                                     str_output,
                                     example_content,
                                     import_content)

        if not os.path.exists(dirpath):
            os.makedirs(dirpath)
        with io.open('{}/{}{}'.format(dirpath, resource_name, extention),
                     'w', encoding='utf8') as f:
            f.write(content_file)
                
    for k, v in navbar_data_source.items():
        navbar_file += """
          <li>
            <a href="{}">{}</a>
          </li>
""".format(v, k)
    
    navbar_file += """
        </ul>
      </li>
      <li>
        <a href="#">Resources</a>
        <ul class="nav">
"""
    for k, v in navbar_resource.items():
        navbar_file += """
          <li>
            <a href="{}">{}</a>
          </li>
""".format(v, k)

    navbar_file += """
        </ul>
      </li>
    </ul>

    <%= partial("layouts/otherdocs", :locals => { :skip => "Terraform Enterprise" }) %>
  <% end %>
  <%= yield %>
<% end %>
"""
    if not os.path.exists(navbar_dirpath):
        os.makedirs(navbar_dirpath)
    if not os.path.exists(index_dirpath):
        os.makedirs(index_dirpath)
    with io.open('{}/outscale.erb'.format(navbar_dirpath),
                 'w', encoding='utf-8') as f:
        f.write(navbar_file)

if __name__ == '__main__':
    main()
