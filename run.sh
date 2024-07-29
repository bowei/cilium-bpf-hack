#!/bin/bash

make
./cfg -mode rawcg -in out.ll -an annotations.txt -start $1 > /tmp/out.gv
dot -Tpdf /tmp/out.gv -o out.pdf
