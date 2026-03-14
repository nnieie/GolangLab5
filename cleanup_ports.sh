#!/usr/bin/env bash

set -euo pipefail

PORTS=(
  8888
  11393
  11394
  11395
  11396
  11397
  6060
  6061
  6062
  6063
  6064
  6065
)

list_listener_pids() {
  local port pid
  declare -A seen=()

  for port in "${PORTS[@]}"; do
    while IFS= read -r pid; do
      [[ -n "${pid}" ]] || continue
      if [[ -z "${seen[$pid]+x}" ]]; then
        seen["$pid"]=1
        printf '%s\n' "$pid"
      fi
    done < <(lsof -tiTCP:"$port" -sTCP:LISTEN 2>/dev/null || true)
  done
}

print_ports() {
  printf '%s ' "${PORTS[@]}"
}

main() {
  local -a pids=()
  local -a remaining=()

  mapfile -t pids < <(list_listener_pids)
  if (( ${#pids[@]} == 0 )); then
    echo "[cleanup] No listeners found on dev ports: $(print_ports)"
    return 0
  fi

  echo "[cleanup] Releasing dev ports: $(print_ports)"
  echo "[cleanup] Sending SIGTERM to PIDs: ${pids[*]}"
  kill -TERM "${pids[@]}" 2>/dev/null || true
  sleep 1

  mapfile -t remaining < <(list_listener_pids)
  if (( ${#remaining[@]} == 0 )); then
    echo "[cleanup] Ports released."
    return 0
  fi

  echo "[cleanup] Still busy after SIGTERM, sending SIGKILL to PIDs: ${remaining[*]}"
  kill -KILL "${remaining[@]}" 2>/dev/null || true
  sleep 1
  echo "[cleanup] Forced cleanup finished."
}

main "$@"
