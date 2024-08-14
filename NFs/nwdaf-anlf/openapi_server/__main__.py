#!/usr/bin/env python3

import connexion

from openapi_server import encoder

# from multiprocessing import Process
import subprocess


def main():
    app = connexion.App(__name__, specification_dir="./openapi/")
    app.app.json_encoder = encoder.JSONEncoder
    app.add_api("openapi.yaml", arguments={"title": "Nnwdaf_AnalyticsInfo"}, pythonic_params=True)

#    app = connexion.App(__name__, specification_dir="/home/free5gc/free5gc/NFs/nwdaf_mtlf/openapi_server/openapi/")
#    app.app.json_encoder = encoder.JSONEncoder
#    app.add_api("openapi.yaml", arguments={"title": "Nnwdaf_MLModelProvision"}, pythonic_params=True)

    # cmd = "/usr/local/go/bin/go run cmd/main.go --nwdafcfg config/nwdafcfg-anlf.yaml" # Example
    # cmd = "../nwdaf-anlf --nwdafcfg ../config/nwdafcfg-anlf.yaml" # Prod
    cmd = "./tmp/main --nwdafcfg ./config/nwdafcfg-anlf.yaml" # Dev

    subprocess.run([cmd], shell=True)
    app.run(port=8080)


if __name__ == "__main__":
    main()

