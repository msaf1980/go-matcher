#!/usr/bin/env python

import argparse
import subprocess
import re
import os
import sys


def parse_cmdline():
    parser = argparse.ArgumentParser(
        description='Detect source changes and run benchmark for changed part (or all, if base package changed)')

    parser.add_argument('-b', '--base', dest='base', action='store', type=str,
                        help='base commit')

    parser.add_argument('-u', '--head', dest='head', action='store', type=str,
                        help='head commit')

    parser.add_argument('-f', '--files', dest='files', action='store', type=str,
                        help='file with changes (relative paths) list')

    parser.add_argument('-c', '--count', dest='count', action='store', type=int, default=6,
                        help='run count')

    return parser.parse_args()


def addDir(file, dirs, baseDirs, sourcesRegexp):
    if file != '':
        dir = os.path.dirname(file)
        if os.path.isdir(dir):        
            for s in sourcesRegexp:
                if s.match(file):
                    dirs.add('./' + dir)
                    for b in baseDirs:
                        if file.startswith(b):
                            # this will trigger all tests
                            return True

    return False


def main():
    args = parse_cmdline()

    # regexp for detect files, which will trigger run benchmark
    sources = [
        re.compile('^.*\.go$')
    ]
    # dirs, which will be trigger all benchmark (not changed dir)
    base = ['pkg/items1/']
    # dirs for all benchmark
    baseTests = ['./...']

    dirs = set()

    isBase = False

    if args.files is None:
        if args.base is None:
            sys.exit("base commit not set")
        if args.head is None:
            sys.exit("head commit not set")

        p = subprocess.run(
            ['git', 'diff', '--name-only', args.base, args.head],
            stdout=subprocess.PIPE,
        )
        for file in p.stdout.decode().split("\n"):
            if addDir(file, dirs, base, sources):
                isBase = True
    else:
        with open(args.files) as f:
            for file in f:
                if addDir(file, dirs, base, sources):
                    isBase = True

    # prerare test dirs list
    tests = []
    if isBase:
        tests = baseTests
    elif len(dirs) > 0:
        for d in dirs:
            tests.append(d)
        tests.sort()

    if len(tests) > 0:
        command = "set -euo pipefail; " \
            "for i in {1..%d}; do"\
            "  echo STEP ${i} ; go test -benchmem -run=^$ -bench '^Benchmark' %s ; "\
            "done" % (args.count, ' '.join(tests))

        sys.stderr.write(command+"\n")
        p = subprocess.Popen([command], shell=True)
        try:
            p.wait()
        except KeyboardInterrupt:
            try:
                p.terminate()
            except OSError:
                pass
            p.wait()


if __name__ == "__main__":
    main()
