#!/bin/sh
for f in examples/*.psph
do
  ./bin/persephone -i $f <<<EOF test EOF
done

for f in examples/conversions/*.psph
do
  ./bin/persephone -i $f
done