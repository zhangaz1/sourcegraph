#!/bin/bash

set -eu

node_modules/.bin/esbuild client/shared/src/api/extension/main.worker.ts --bundle --outfile=ui/assets/worker.js --loader:.yaml=text --define:global=window

node_modules/.bin/esbuild client/web/src/enterprise/main.tsx --bundle '--define:process.env.NODE_ENV="development"' --outfile=ui/assets/out.js --loader:.yaml=text --loader:.scss=text --define:global=window
