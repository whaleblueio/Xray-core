#!/bin/sh
# vim:sw=4:ts=4:et

set -e

if /usr/bin/find "/docker-entrypoint.d/" -mindepth 1 -maxdepth 1 -type f -print -quit 2>/dev/null | read v; then
        echo >&1 "$0: /docker-entrypoint.d/ is not empty, will attempt to perform configuration"

  find "/docker-entrypoint.d/" -follow -type f -print | sort -n | while read -r f; do
      case "$f" in
          *.sh)
              if [ -x "$f" ]; then
                  echo >&1 "$0: Launching $f";
                  "$f"
              else
                  # warn on shell scripts without exec bit
                  echo >&1 "$0: Ignoring $f, not executable";
              fi
              ;;
          *) echo >&1 "$0: Ignoring $f";;
      esac
  done
else
        echo >&1 "$0: No files found in /docker-entrypoint.d/, skipping configuration"
    fi
exec "$@"
