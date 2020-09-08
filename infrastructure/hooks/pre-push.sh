#!/bin/sh

if [ -f ./.git/hooks/old-pre-push ]; then source ./.git/hooks/old-pre-push; fi

make pre_push
