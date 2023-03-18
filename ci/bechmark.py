#!/usr/bin/env python

import argparse
import subprocess
import re
import os.path


def parse_cmdline():
    parser = argparse.ArgumentParser(
        description='Detect source changes and run benchmark for changed part (or all, if base package changed)')

    parser.add_argument('-f', '--from', dest='base', action='store', type=str, required=True,
                        help='base commit')

    parser.add_argument('-u', '--until', dest='head', action='store', type=str, required=True,
                        help='head commit')

    parser.add_argument('-c', '--count', dest='count', action='store', type=int, default=6,
                        help='run count')

    return parser.parse_args()


def main():
    args = parse_cmdline()

    # regexp for detect files, which will trigger run benchmark
    sources = [
        re.compile('^.*\.go$')
    ]
    # dirs, which will be trigger all benchmark (not changed dir)
    base = ['pkg/items/']
    # dirs for all benchmark
    baseTests = ['./...']

    dirs = set()

    isBase = False

    p = subprocess.run(
        ['git', 'diff', '--name-only', args.base, args.head],
        stdout=subprocess.PIPE,
    )
    for file in p.stdout.decode().split("\n"):
        if file == '':
            continue
        for s in sources:
            if not isBase:
                for b in base:
                    if file.startswith(b):
                        isBase = True
                if s.match(file):
                    dirs.add('./' + os.path.dirname(file))

    tests = []
    if isBase:
        tests = baseTests
    else:
        for d in dirs:
            tests.append(d)
        tests.sort()

    if len(tests) > 0:
        command = "set -euo pipefail; " \
            "for i in {1..%d}; do"\
            "  echo STEP ${i} ; go test -benchmem -run=^$ -bench '^Benchmark' %s ; "\
            "done" % (args.count, ' '.join(tests))

        print(command)
        subprocess.run([command], shell=True)


if __name__ == "__main__":
    main()
