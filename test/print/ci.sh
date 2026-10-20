#!/bin/bash
# Prints documents of testdata with bdf print on printers that are not
# bdf's own: ippeveprinter (CUPS's IPP Everywhere printer) directly, as
# PDF and as PWG Raster, a CUPS queue in front of it (an IPP Everywhere
# queue, made with lpadmin -m everywhere), and bdf printd's spool in
# front of it across a restart, and checks what the printer received.
# The CI job "print" runs it.
#
# Needs ippeveprinter (cups-ipp-utils), cupsd and lpadmin (cups) with
# its filters, a running Avahi (ippeveprinter will not start without
# DNS-SD) and IPv6 on the loopback (ippeveprinter listens on both), and
# sudo for lpadmin. BDF names the bdf command (default: go run).
set -euo pipefail
root=$(cd "$(dirname "$0")/../.." && pwd)
cd "$root"
bdf=${BDF:-}
if [ -z "$bdf" ]; then
  bdf=$(mktemp -d)/bdf
  go build -o "$bdf" ./cmd/bdf
fi
spool=$(mktemp -d)
chmod 777 "$spool"
fonts=(-font-dir converter/pptx/testdata/fonts -no-system-fonts)
port=8631
url="ipp://localhost:$port/ipp/print"

ippeveprinter -2 -p "$port" -n localhost -k -d "$spool" -s 120 \
  -f application/pdf,image/pwg-raster,image/jpeg "BDF Test" >"$spool/ippeveprinter.log" 2>&1 &
pid=$!
trap 'kill $pid ${printd:-} 2>/dev/null || true; sudo lpadmin -x bdftest 2>/dev/null || true' EXIT
for _ in $(seq 30); do
  if "$bdf" printers -url "$url" >"$spool/printer.txt" 2>/dev/null; then break; fi
  sleep 1
done
cat "$spool/printer.txt"
grep -q '^formats: .*application/pdf' "$spool/printer.txt"
grep -q '^formats: .*image/pwg-raster' "$spool/printer.txt"
grep -q '^sides: .*two-sided-long-edge' "$spool/printer.txt"

# what ippeveprinter kept of the jobs: the documents whose first bytes
# are these
received() {
  local n=0 f
  for f in "$spool"/*; do
    [ -f "$f" ] || continue
    case "$f" in *.log|*.txt) continue ;; esac
    if [ "$(head -c ${#1} "$f")" = "$1" ]; then n=$((n + 1)); fi
  done
  echo "$n"
}

# PDF, two-sided, two copies, A4
"$bdf" print -printer "$url" "${fonts[@]}" -wait 5m -sides long -copies 2 -media a4 testdata/form/rows.bdf
test "$(received '%PDF-')" -ge 1

# PWG Raster in black and white at 300 dpi, two-sided: the back of the
# sheet as the printer's pwg-raster-document-sheet-back asks
"$bdf" print -printer "$url" "${fonts[@]}" -wait 5m -format pwg-raster -mono -bilevel -resolution 300 -sides long testdata/form/rows.bdf
test "$(received 'RaS2')" -ge 1

# a page left out of what the printer does not support: -strict fails
# with exit status 4
set +e
"$bdf" print -printer "$url" "${fonts[@]}" -strict -source tray-99 testdata/form/rows.bdf
status=$?
set -e
test "$status" -eq 4

# a CUPS queue of the printer, as Linux and macOS have them
sudo lpadmin -p bdftest -E -v "$url" -m everywhere
"$bdf" printers -cups local | grep -q '^bdftest'
before=$(received '%PDF-')
"$bdf" print -printer ipp://localhost:631/printers/bdftest "${fonts[@]}" -wait 5m testdata/form/pptx-fixed.bdf
# CUPS sends the job on to ippeveprinter
for _ in $(seq 60); do
  if [ "$(received '%PDF-')" -gt "$before" ] || [ "$(received 'RaS2')" -gt 1 ]; then break; fi
  sleep 1
done
ls -la "$spool"
test "$(received '%PDF-')" -gt "$before" || test "$(received 'RaS2')" -gt 1

# bdf printd: two jobs, one canceled; the server stops right after the
# other is submitted and starts again; the printer has it once
work=$(mktemp -d)
server=http://127.0.0.1:8640
cat >"$work/printd.yaml" <<YAML
listen: 127.0.0.1:8640
printers: { eve: { url: "$url" } }
queues: { office: { printers: [eve], defaults: { sides: two-sided-long-edge } } }
fonts: { dirs: ["$root/converter/pptx/testdata/fonts"], no-system: true }
YAML
"$bdf" printd -config "$work/printd.yaml" -check
start_printd() {
  "$bdf" printd -config "$work/printd.yaml" -log text 2>>"$work/printd.log" &
  printd=$!
  for _ in $(seq 30); do
    if curl -sf "$server/healthz" >/dev/null; then return; fi
    sleep 1
  done
  cat "$work/printd.log"
  return 1
}
start_printd
"$bdf" printers -server "$server"
before=$(received '%PDF-')
held=$("$bdf" print -server "$server" -queue office -q -hold indefinite testdata/form/rows.bdf)
"$bdf" jobs -server "$server" cancel "$held"
job=$("$bdf" print -server "$server" -queue office -q testdata/form/pptx-fixed.bdf)
kill -INT "$printd"
wait "$printd"
start_printd
for _ in $(seq 120); do
  if "$bdf" jobs -server "$server" show "$job" | grep -q '^state: *completed'; then break; fi
  sleep 1
done
"$bdf" jobs -server "$server" -all
"$bdf" jobs -server "$server" show "$job" | grep -q '^state: *completed'
test "$(received '%PDF-')" -eq $((before + 1))
"$bdf" jobs -server "$server" accounting -csv | tee "$work/accounting.csv"
test "$(wc -l <"$work/accounting.csv")" -eq 3
kill -INT "$printd"
wait "$printd"
printd=
"$bdf" printd -config "$work/printd.yaml" -compact
cat "$work/printd.log"
echo "print: ok"
