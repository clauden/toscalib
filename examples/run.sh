#!/bin/bash


EXECUTABLE="./my_parser.go"
YAML="./ecomp_20161215.yml"
OUTPUT=/tmp/$$

executable="$EXECUTABLE"
output="$OUTPUT"
dryrun=0
yaml=""

usage() {
  echo "Usage: %0 [-x executable] [-y yaml-file] [-o output-file] -d"
  echo "       d means dry-run"
  echo
}

options=":x:y:o:d"
while getopts $options opt ; do
  case $opt in
    x)
      executable="$OPTARG"
      echo "Exec: $executable" >&2
      ;;

    y)
      yaml="$OPTARG"
      echo "Yaml: $yaml" >&2
      ;;

    o)
      output="$OPTARG"
      echo "Output: $output" >&2
      ;;

    d)
      dryrun=1
      echo "Dry run" >&2
      ;;

    \?)
      echo "Unknown option -$OPTARG" >&2
      usage
      exit 1
      ;;

    :)
      echo "Missing argument for -$OPTARG" >&2
      usage
      exit 1
      ;;

    *)
      echo "Unimplemented -$OPTARG" >&2
      usage 
      exit 1
      ;;
  esac
done
shift $((OPTIND-1))

if [[ -z "$yaml" ]]; then
  if [[ "$1" ]]; then
    yaml="$1"
  else
    yaml="$YAML"
  fi
fi

echo "go run $executable < $yaml > $output.dot"

go run $executable < $yaml > $output.dot && dot -o $output.png -Tpng $output.dot && open $output.png

