#!/usr/bin/env sh
set -eu

BIN="${BIN:-./bin/yunxiao-mcp-server}"
PORT="${PORT:-39393}"
HEALTH_URL="http://127.0.0.1:${PORT}/healthz"
LOG_FILE="${TMPDIR:-/tmp}/yunxiao-mcp-server-smoke.$$.log"

command -v curl >/dev/null 2>&1 || {
	echo "smoke failed: curl is required" >&2
	exit 1
}
command -v nc >/dev/null 2>&1 || {
	echo "smoke failed: nc is required" >&2
	exit 1
}

if nc -z 127.0.0.1 "${PORT}" >/dev/null 2>&1; then
	echo "smoke failed: port ${PORT} is already in use" >&2
	exit 1
fi

"${BIN}" version | grep "yunxiao-mcp-server version=" >/dev/null
"${BIN}" --help >/dev/null

"${BIN}" --port "${PORT}" --access-token smoke-token --log-level disabled >"${LOG_FILE}" 2>&1 &
SERVER_PID=$!

cleanup() {
	kill "${SERVER_PID}" >/dev/null 2>&1 || true
	wait "${SERVER_PID}" >/dev/null 2>&1 || true
	rm -f "${LOG_FILE}"
}
trap cleanup EXIT INT TERM

i=0
while [ "${i}" -lt 50 ]; do
	i=$((i + 1))
	if ! kill -0 "${SERVER_PID}" >/dev/null 2>&1; then
		cat "${LOG_FILE}" >&2 || true
		exit 1
	fi
	if response="$(curl -fsS "${HEALTH_URL}" 2>/dev/null)" && [ "${response}" = "healthy" ]; then
		if ! kill -0 "${SERVER_PID}" >/dev/null 2>&1; then
			cat "${LOG_FILE}" >&2 || true
			exit 1
		fi
		exit 0
	fi
	sleep 0.1
done

cat "${LOG_FILE}" >&2 || true
echo "smoke failed: ${HEALTH_URL} did not become healthy" >&2
exit 1
