#!/bin/sh
set -euo

CALICO_VERSION=$1
CALICO_VERSION_TRIMMED=$(echo "$CALICO_VERSION" | sed 's/v//') # trim leading v
TMPDIR=$(mktemp -d)
# for example, https://github.com/projectcalico/calico/tree/v3.16.0/_includes/charts/calico/crds/kdd
wget https://github.com/projectcalico/calico/archive/refs/tags/"$CALICO_VERSION".zip -O "$TMPDIR"/calico.zip
unzip -j "$TMPDIR"/calico.zip calico-"$CALICO_VERSION_TRIMMED"/_includes/charts/calico/crds/kdd/\*.yaml -d crds
