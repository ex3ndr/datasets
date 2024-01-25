import fs from 'fs';
import crypto from 'crypto';

// Read all files
let path = '/Volumes/shared/datasets/common-voice/';
let files = fs.readdirSync(path);
files = files.filter((v) => v.startsWith('cv-corpus-16.0-2023-12-06-') && v.endsWith('.tar.gz'));
let langs = files.map((v) => v.substring('cv-corpus-16.0-2023-12-06-'.length, v.length - 7));

// Hash all files
function calculateHashes(filename, algorithms) {
    return new Promise((resolve, reject) => {
        const hashes = algorithms.map((v) => crypto.createHash(v));
        const stream = fs.createReadStream(filename);
        stream.on('error', err => reject(err));
        stream.on('data', chunk => hashes.forEach((v) => v.update(chunk)));
        stream.on('end', () => resolve(hashes.map((v) => v.digest('hex'))));
    });
}

// Hash all files
let hashes = {};
for (let f of files) {
    console.log('Hashing ' + f);
    let res = await calculateHashes(path + f, ['md5', 'sha1', 'sha256']);
    hashes[f] = res;
}

// Write collection
let endpoint = 'https://shared.korshakov.com/datasets/common-voice/';
let template = `# Description
id: {id}
name: "{name}"
url: "https://commonvoice.mozilla.org"
description: "Mozilla Common Voice is an initiative to help teach machines how real people speak."
license: CC BY 4.0

# Dataset
dataset:
  url: "{dataset_url}"
  hashes:
    sha1: {sha1}
    sha256: {sha256}
    md5: {md5}
`;
for (let l of langs) {

    let fname = 'cv-corpus-16.0-2023-12-06-' + l + '.tar.gz';
    let h = hashes[fname];
    let dataset_url = endpoint + fname;
    let id = 'common-voice-16.0-' + l;
    let name = 'Common Voice 16.0 (' + l + ')';
    let sha1 = h[1];
    let sha256 = h[2];
    let md5 = h[0];
    let content = template
        .replace('{id}', id)
        .replace('{name}', name)
        .replace('{dataset_url}', dataset_url)
        .replace('{sha1}', sha1)
        .replace('{sha256}', sha256).replace('{md5}', md5);

    fs.writeFileSync(__dirname + '/../collection/common-voice-16.0-' + l + '.yaml', content);
}