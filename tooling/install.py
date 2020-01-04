import argparse
import subprocess
import sys
import os

sys.path.insert(0, os.path.abspath(
    os.path.join(os.path.dirname(__file__), '..')))

import yap

parser = argparse.ArgumentParser()
parser.add_argument("-o", "--output", help="path for installation",
                    default="/usr/local/bin/")
args = parser.parse_args()
# todo read env
install_file = args.output

tag = "unknown"
try:
    tag = subprocess.run(["git", "describe", "--tags"],
                         check=True,
                         capture_output=True,
                         encoding="utf-8").stdout.strip()
except Exception as e:
    print("Loading git tag failed with: {}".format(e))

commit = ""
try:
    commit = subprocess.run(["git", "rev-parse", "--short", "HEAD"],
                            check=True,
                            capture_output=True,
                            encoding="utf-8").stdout.strip()
except Exception as e:
    print("Loading git commit failed with: {}".format(e))

p = yap.Pipeline(settings={
    "state": {
        "type": "json"
    }
})
p.add(yap.DictBlock({
    "name": "test",
    "exec": "go test -v -count=1 ./...",
    "in_files": ["**/*.go"]
}))
build_cmd = 'go build -ldflags "' \
            '-X github.com/tivvit/yap/cmd.GitTag={} ' \
            '-X github.com/tivvit/yap/cmd.GitCommit={}" ' \
            '.'.format(tag, commit)
p.add(yap.DictBlock({
    "name": "build",
    "exec": build_cmd,
    "in_files": ["**/*.go"]
}))
p.add(yap.DictBlock({
    "name": "install",
    "exec": "cp yap {}".format(install_file),
    "deps": ["build"],
    "in_files": ["**/*.go"]
}))
print(p)
