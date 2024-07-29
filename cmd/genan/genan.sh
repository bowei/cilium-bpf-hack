#!/bin/bash

find "$1" \
    \( -name \*.c -o -name \*.h \) \
    -exec ./genan -strip "$1/" {} \;