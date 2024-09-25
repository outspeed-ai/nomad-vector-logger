"""
This script is used in user-data-client.sh script to update the values of config.toml file which the nomad-vector-logger will use to generate the vector config file.
"""

import sys
import tomli
import tomli_w

override = {"loki": {}}

for arg in sys.argv[1:]:
    pair = arg.split("=")
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


# Read the config.sample.toml file
with open("config.sample.toml", "rb") as f:
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

with open("config.toml", "wb") as f:
    tomli_w.dump(config, f)
    print("👍 config.toml generated")

# this generated config.toml will then be used by nomad-vector-logger program
