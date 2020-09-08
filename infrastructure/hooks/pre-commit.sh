#!/bin/sh

if [ -f ./.git/hooks/old-pre-commit ]; then source ./.git/hooks/old-pre-commit; fi

make pre_commit
