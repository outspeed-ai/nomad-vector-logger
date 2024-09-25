"""
This script is used in user-data-client.sh script to update the values of nomad-vector-logger.toml.tpl
file which the nomad-vector-logger will use to generate the vector config file.
"""

import os
import sys
import tomli
import tomli_w

override = {"loki": {}}

for arg in sys.argv[1:]:
    pair = arg.split("=", maxsplit=1)
    if len(pair) != 2:
        print("error: invalid argument", arg)
        exit(1)

    key, value = pair
    key = key.strip().lower()
    value = value.strip()
    if key.startswith("loki_"):
        override["loki"][key[5:]] = value
    else:
        override[key] = value

CONFIG_FILE = "./dev/nomad-vector-logger.toml.tpl"

# Read the config.sample.toml file
with open(CONFIG_FILE, "rb") as f:
    config = tomli.load(f)

    # update the config with the values from the command line arguments
    # update loki config
    for key, value in override["loki"].items():
        config["app"]["loki"][key] = value

    # remove loki key from override since it's already been added to the config
    del override["loki"]

    # update the rest of the config
    for key, value in override.items():
        config["app"][key] = value

# backup the original config file
os.rename(CONFIG_FILE, f"{CONFIG_FILE}.bak")

# write to the ori
with open(CONFIG_FILE, "w") as f:
    """
    # the following lines need to be like this only
    
    vector_config_dir = "{{ env "NOMAD_ALLOC_DIR" }}/vector_gen_configs"
    extra_templates_dir = "{{ env "NOMAD_TASK_DIR" }}/static/"


    but the generated config.toml will have the quotations escaped
    so we need to replace them with double quotes
    """

    config_str = tomli_w.dumps(config).replace('\\"', '"')
    f.write(config_str)
    print("👍 config.toml generated")

# this generated config.toml will then be used by nomad-vector-logger program
