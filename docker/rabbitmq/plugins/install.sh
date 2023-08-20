#!/bin/bash
# shellcheck disable=SC2207
pls=($(ls))
containerName=rabbitmq
aln=${#pls[@]}
ez=.ez
for ((i = 0; i < aln; i++)); do
  va=${pls[i]}
  if [[ $va != *$ez ]]; then
    continue
  fi
  echo $va
  docker cp $va $containerName:/opt/rabbitmq/plugins/
done
