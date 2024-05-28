#!/usr/bin/env bash
image='api-units'

printerror () {
  echo "=========================================================================================="
  echo "=";
  echo "= FAILURE: $1";
  echo "=";
  echo "=========================================================================================="
  exit 1;
}

usage() {
    echo "Usage: $0 [-t <tag name>] -n <test|preprod|prod>" 1>&2;
    echo
    echo "Ex.: $0 -t develop -n md-api-test" 1>&2;
    echo
    exit 1;
}

while getopts ":t:n:" option; do
    case "${option}" in
        t)
            tagname=${OPTARG}
            ;;
        n)
            namespace=${OPTARG}
            ;;
        *)
            usage
            ;;
    esac
done
shift $((OPTIND-1))

if [ -z "${namespace}" ]; then
    usage
fi

env=`echo $namespace | grep -oP "[a-z0-9]+\-\K(test|preprod|prod)"`

# Checkout if tagname specified, otherwise use current code state
if [ ! -z "${tagname}" ]; then
  # Checkout code on tag
  git checkout $tagname
  if [ $? -ne 0 ]; then printerror "Couldn't checkout code at $tagname"; fi
fi

# Set default value to develop if not provided
if [ -z "${tagname}" ]; then
    tagname='develop'
fi

# Check that you are connected to OS and select namespace
oc project $namespace
if [ $? -ne 0 ]; then printerror 'oc failure, make sure it is installed, that you are logged in, and the project exists of you are granted on it'; fi

# Build image
docker build -t $image:$tagname .
if [ $? -ne 0 ]; then printerror 'could not build image'; fi

# Tag image with imagetag
docker tag $image:$tagname os-docker-registry.epfl.ch/$namespace/$image:$tagname
# Push image
docker login os-docker-registry.epfl.ch
docker push os-docker-registry.epfl.ch/$namespace/$image:$tagname

if [ $? -ne 0 ]; then printerror "could not push image to registry, check you are logged in on the right namespace ($namespace)"; fi

oc apply -f ../md-api-infra/$env/deployment.yml
