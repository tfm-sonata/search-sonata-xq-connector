#!/usr/bin/env bash
rm -rf ./.serverless

export version=$(cat ./version.txt)
export accountNo=906404173569
export profile=tfm-develop
export role=afklm_developer

export stage=v1_1_0
export partner=TFM
export name=search-sonata-connector
export service_name=SEARCH-SONATA-CONNECTOR
export stack=search-sonata-connector-CFStack
export bucket=tfm-serverless-dev
export environment=tfm

export sg1=sg-0630d13ea4b138158             #SSH
export sg2=sg-2ed9fd4d                      #default
export sn1=subnet-b76fbddd         #app-subnet-1
export sn2=subnet-d12028ac

export config_provider=common-config-service
export queue_in=search-flight-xq

#make
serverless deploy --verbose
echo Script finished on `date`