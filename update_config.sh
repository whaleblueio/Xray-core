#!/bin/sh
set -e
ME=$(basename $0)
auto_envsubst() {
  local template_dir="${ENVSUBST_TEMPLATE_DIR:-/etc/trojan-go}"
  local suffix="${ENVSUBST_TEMPLATE_SUFFIX:-.template}"
  local output_dir="${ENVSUBST_OUTPUT_DIR:-/etc/trojan-go}"

  local template defined_envs relative_path output_path subdir
  defined_envs=$(printf '${%s} ' $(env | cut -d= -f1))
  [ -d "$template_dir" ] || return 0
  if [ ! -w "$output_dir" ]; then
    echo >&1 "$ME: ERROR: $template_dir exists, but $output_dir is not writable"
    return 0
  fi

  find "$template_dir" -follow -type f -name "*$suffix" -print | while read -r template; do
    relative_path="${template#$template_dir/}"
    output_path="$output_dir/${relative_path%$suffix}"
    subdir=$(dirname "$relative_path")
    # create a subdirectory where the template file exists
    mkdir -p "$output_dir/$subdir"
    echo >&1 "$ME: Running envsubst on $template to $output_path"
    envsubst "$defined_envs" < "$template" > "$output_path"

  done
}

auto_envsubst

exit 0


