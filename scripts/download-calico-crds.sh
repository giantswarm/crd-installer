#!/bin/sh
set -euo

# This script isn't used directly in this repo. Instead, it is called during retagging (e.g.
# https://github.com/giantswarm/retagger/blob/c965cc81bde620f56e3a71cbea7f1b9c95b37544/images.yaml#L599)

CALICO_VERSION=$1
CALICO_VERSION_TRIMMED=$(echo "$CALICO_VERSION" | sed 's/v//') # trim leading v
TMPDIR=$(mktemp -d)
# Example final URL: https://github.com/projectcalico/calico/tree/v3.16.0/_includes/charts/calico/crds/kdd
wget https://github.com/projectcalico/calico/archive/refs/tags/"$CALICO_VERSION".zip -O "$TMPDIR"/calico.zip

unzip -j "$TMPDIR"/calico.zip calico-"$CALICO_VERSION_TRIMMED"/libcalico-go/config/crd/\*.yaml -d crds

count="$(ls -1 crds|wc -l)"

echo "Found $count crds"

if [ $count -eq 0 ]
then
  echo "Can't find any CRD in the calico zip file for version $CALICO_VERSION"
  exit 1
fi
