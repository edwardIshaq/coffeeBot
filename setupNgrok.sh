#!/bin/sh

set -o pipefail
set -x

ngrok http -subdomain=goPlatform  8080
