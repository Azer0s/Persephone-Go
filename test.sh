#!/bin/bash

for f in examples/*.psph
do
  if [[ "$f" != "examples/input.psph" ]]; then
    ./bin/persephone -i $f
  fi
done

for f in examples/conversions/*.psph
do
  ./bin/persephone -i $f
done
