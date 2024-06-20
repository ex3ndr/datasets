set -e
NAME="$1"

# Pack dataset
mv dataset $NAME
COPYFILE_DISABLE=1 tar -cf $NAME.tar $NAME
mv $NAME dataset

# Checksums without filename
sha1=`shasum -a 1 $NAME.tar | cut -d ' ' -f 1`
sha256=`shasum -a 256 $NAME.tar | cut -d ' ' -f 1`
md5=`md5 $NAME.tar | cut -d ' ' -f 4`

echo "sha1: $sha1"
echo "sha256: $sha256"
echo "md5: $md5"