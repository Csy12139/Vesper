#!/bin/bash

# Exit on error
set -xe

./build/mn ../MN/mn_conf_template.json

./build-test/dn ../DN/dn_conf_template.json