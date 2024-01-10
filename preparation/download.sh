set -e

# Reset directory
rm -fr ./dataset
mkdir ./dataset

# Download dataset
wget "$1" -O dataset.tar.gz

# Unpack dataset
tar -xvf ./dataset.tar.gz -C ./dataset --strip-components=1

# Delete tarball
rm ./dataset.tar.gz