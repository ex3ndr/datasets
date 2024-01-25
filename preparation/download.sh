set -e

# Reset directory
rm -fr ./dataset
mkdir ./dataset

# Download dataset
wget "$1" -O dataset.tar.gz
# wget "$1" -O dataset.zip

# Unpack dataset
tar -xvf ./dataset.tar.gz -C ./dataset --strip-components=1
# bsdtar -xvf ./dataset.zip -C ./dataset --strip-components=1

# Delete tarball
rm ./dataset.tar.gz
# rm ./dataset.zip