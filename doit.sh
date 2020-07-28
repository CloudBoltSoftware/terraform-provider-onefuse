#!/bin/bash

# clean, build, setup environment, init and apply the Terraform provider

rm -f terraform.tfstate
rm -rf /tmp/tf-state*
rm /tmp/terraform-log
echo 'building go plugin...'
go build -o terraform-provider-onefuse
echo 'setting environment...'
source ./setenv.sh
echo 'running terraform...'
terraform init
terraform plan
terraform apply

