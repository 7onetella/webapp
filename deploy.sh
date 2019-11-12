#!/bin/sh
set +x

"${GOPATH}"/bin/gox -osarch="linux/amd64"

BUILD_ID=$1

chmod +x webapp_linux_amd64

BINARY=webapp_${BUILD_ID}
mv webapp_linux_amd64 ${BINARY}    

echo uploading to file server
TAR_FILE=webapp_${BUILD_ID}.tar.gz
tar czvf ${TAR_FILE} ${BINARY}

scp ${TAR_FILE} jenkins@vagrant.7onetella.net:/mnt/uploads
rm ./${BINARY}
rm ./${TAR_FILE}

echo scheduling deployment
cat ./webapp.nomad | sed 's|BUILD_ID|'"${BUILD_ID}"'|g' > webapp.nomad.${BUILD_ID}
export NOMAD_ADDR=http://nomad.7onetella.net:4646
nomad job run ./webapp.nomad.${BUILD_ID}
rm ./webapp.nomad.${BUILD_ID}
