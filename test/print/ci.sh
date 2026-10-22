#!/bin/bash
# Prints documents of testdata with bdf print on fake printers
# (print/fakeprinter, run as a farm by tools/fakeprinter) and on what CUPS
# makes of them, and checks what the printers received; Ghostscript draws
# the PostScript. The fakes: an IPP Everywhere printer, as PDF and as PWG
# Raster; one of PostScript only; a CUPS queue in front of the first (an
# IPP Everywhere queue, made with lpadmin -m everywhere, as CUPS makes it
# of the fake's attributes); a printer of the raw port 9100; then bdf
# printd's spool in front of them across a restart. Then CUPS's own
# clients on the fakes (raw queues of 9100 and LPD, lp). Then bdf printd
# as an IPP printer: ipptool's checks of IPP Everywhere, jobs of ipptool,
# its DNS-SD records, and a CUPS queue made of it. Last, the languages of
# labels and receipts: ZPL over 9100, and ESC/POS through a raw queue of
# CUPS.
# (The fakes themselves are checked against ipptool and ippeveprinter by
# go test ./print/fakeprinter/, which the CI job runs first.)
# The CI job "print" runs it.
#
# Needs cupsd and lpadmin (cups) with its filters, ipptool
# (cups-ipp-utils), a running Avahi (for the DNS-SD of bdf printd), sudo
# for lpadmin and Ghostscript (gs). BDF names the bdf command (default:
# go run).
set -euo pipefail
# (a command piped into grep writes to a grep that reads all of it: grep -q
# quits at its first match, and the command would die of SIGPIPE, which
# pipefail makes the script's failure)
root=$(cd "$(dirname "$0")/../.." && pwd)
cd "$root"
bdf=${BDF:-}
if [ -z "$bdf" ]; then
  bdf=$(mktemp -d)/bdf
  go build -o "$bdf" ./cmd/bdf
fi
fonts=(-font-dir converter/pptx/testdata/fonts -no-system-fonts)
work=$(mktemp -d)
# the farm keeps what each printer received in a directory of its own
store=$(mktemp -d)
chmod 777 "$store"
port=8631
url="ipp://127.0.0.1:$port/ipp/print"
spool="$store/eve"
psport=8632
psurl="ipp://127.0.0.1:$psport/ipp/print"
psspool="$store/ps"
rawport=9101
rawdir="$store/raw"
lpdport=8515
mkdir -p "$spool" "$psspool" "$rawdir"
cat >"$work/farm.yaml" <<YAML
listen: 127.0.0.1
web: 127.0.0.1:8650
store: $store
printers:
  - { name: eve, preset: office-mfp, port: $port, formats: [application/pdf, image/pwg-raster, image/jpeg, application/octet-stream], processing: 1s }
  - { name: ps, preset: ps-office, port: $psport, formats: [application/postscript, application/octet-stream] }
  - { name: raw, preset: jetdirect-ps, port: $rawport }
  - { name: lpd, preset: lpd-legacy, port: $lpdport, queue: lp }
YAML
farm=$(mktemp -d)/fakeprinter
go build -o "$farm" ./tools/fakeprinter
"$farm" -config "$work/farm.yaml" >"$work/farm.log" 2>&1 &
pid=$!
trap 'kill $pid ${printd:-} 2>/dev/null || true; for q in bdftest bdfserver bdfraw fakesock fakelpd; do sudo lpadmin -x $q 2>/dev/null || true; done' EXIT
# the number of jobs a printer of the farm has
jobs_of() {
  curl -sf "http://127.0.0.1:8650/api/jobs?printer=$1" | grep -c '"printer":' || true
}
for _ in $(seq 30); do
  if "$bdf" printers -url "$url" >"$spool/printer.txt" 2>/dev/null; then break; fi
  sleep 1
done
cat "$spool/printer.txt"
grep -q '^formats: .*application/pdf' "$spool/printer.txt"
grep -q '^formats: .*image/pwg-raster' "$spool/printer.txt"
grep -q '^sides: .*two-sided-long-edge' "$spool/printer.txt"

# what a printer kept of the jobs: the documents whose first bytes are
# these (in the farm's directory of the printer, or another)
received() {
  local dir=$spool
  if [ $# -gt 1 ]; then dir=$2; fi
  local n=0 f
  for f in "$dir"/*; do
    [ -f "$f" ] || continue
    case "$f" in *.log|*.txt) continue ;; esac
    if [ "$(head -c ${#1} "$f")" = "$1" ]; then n=$((n + 1)); fi
  done
  echo "$n"
}

# Ghostscript draws a PostScript file without a word
draw() {
  local out
  out=$(mktemp -d)
  gs -q -dSAFER -dBATCH -dNOPAUSE -sDEVICE=png16m -r72 -sOutputFile="$out/p%d.png" "$1" >"$out/gs.log" 2>&1 || { cat "$out/gs.log"; return 1; }
  if [ -s "$out/gs.log" ]; then cat "$out/gs.log"; return 1; fi
  test -s "$out/p1.png"
}

# PDF, two-sided, two copies, A4
"$bdf" print -printer "$url" "${fonts[@]}" -wait 5m -sides long -copies 2 -media a4 testdata/form/rows.bdf
test "$(received '%PDF-')" -ge 1

# PWG Raster in black and white at 300 dpi, two-sided: the back of the
# sheet as the printer's pwg-raster-document-sheet-back asks
"$bdf" print -printer "$url" "${fonts[@]}" -wait 5m -format pwg-raster -mono -bilevel -resolution 300 -sides long testdata/form/rows.bdf
test "$(received 'RaS2')" -ge 1

# PostScript, to the printer of PostScript only: the transparency of the
# soft masks drawn as pictures, the rest vectors
for _ in $(seq 30); do
  if "$bdf" printers -url "$psurl" >"$psspool/printer.txt" 2>/dev/null; then break; fi
  sleep 1
done
cat "$psspool/printer.txt"
grep -q '^formats: .*application/postscript' "$psspool/printer.txt"
"$bdf" print -printer "$psurl" "${fonts[@]}" -wait 5m -sides long -media a4 testdata/form/rows.bdf
"$bdf" print -printer "$psurl" "${fonts[@]}" -wait 5m testdata/pdf/chrome-masks.bdf
test "$(received '%!PS-Adobe-3.0' "$psspool")" -eq 2
for f in "$psspool"/*; do
  if [ "$(head -c 14 "$f")" = '%!PS-Adobe-3.0' ]; then draw "$f"; fi
done

# the raw port 9100: PostScript, then PDF with a PJL header
"$bdf" print -printer "socket://127.0.0.1:$rawport" "${fonts[@]}" -copies 2 -media a4 testdata/form/rows.bdf
"$bdf" print -printer "socket://127.0.0.1:$rawport?format=pdf" "${fonts[@]}" -sides long testdata/form/rows.bdf
test "$(received '%!PS-Adobe-3.0' "$rawdir")" -eq 1
grep -q 'NumCopies 2' "$rawdir"/1-1.*
draw "$(ls "$rawdir"/1-1.*)"
grep -aq '@PJL ENTER LANGUAGE=PDF' "$rawdir"/2-1.*
grep -aq '@PJL SET DUPLEX=ON' "$rawdir"/2-1.*

# a page left out of what the printer does not support: -strict fails
# with exit status 4
set +e
"$bdf" print -printer "$url" "${fonts[@]}" -strict -source tray-99 testdata/form/rows.bdf
status=$?
set -e
test "$status" -eq 4

# a CUPS queue of the printer, as Linux and macOS have them
sudo lpadmin -p bdftest -E -v "$url" -m everywhere
"$bdf" printers -cups local | grep '^bdftest' >/dev/null
before=$(received '%PDF-')
"$bdf" print -printer ipp://localhost:631/printers/bdftest "${fonts[@]}" -wait 5m testdata/form/pptx-fixed.bdf
# CUPS sends the job on to the fake
for _ in $(seq 60); do
  if [ "$(received '%PDF-')" -gt "$before" ] || [ "$(received 'RaS2')" -gt 1 ]; then break; fi
  sleep 1
done
ls -la "$spool"
test "$(received '%PDF-')" -gt "$before" || test "$(received 'RaS2')" -gt 1

# CUPS's own clients on the fakes: raw queues of the raw port and of
# LPD (cupsd's socket and lpd backends), and lp; the job CUPS sent the
# IPP Everywhere queue above is drawn back by the farm
sudo lpadmin -p fakesock -E -v "socket://127.0.0.1:$rawport" -m raw
sudo lpadmin -p fakelpd -E -v "lpd://127.0.0.1:$lpdport/lp" -m raw
"$bdf" pdf "${fonts[@]}" testdata/form/rows.bdf "$work/rows-raw.pdf"
rawbefore=$(jobs_of raw)
lp -d fakesock "$work/rows-raw.pdf"
lp -d fakelpd "$work/rows-raw.pdf"
for _ in $(seq 60); do
  if [ "$(jobs_of raw)" -gt "$rawbefore" ] && [ "$(jobs_of lpd)" -ge 1 ]; then break; fi
  sleep 1
done
test "$(jobs_of raw)" -gt "$rawbefore"
test "$(jobs_of lpd)" -ge 1
curl -sf "http://127.0.0.1:8650/api/jobs?printer=lpd" | grep '"format": "application/pdf"' >/dev/null
last=$(curl -sf "http://127.0.0.1:8650/api/jobs?printer=eve" | awk '/"id":/ && !n { print; n = 1 }' | tr -dc '0-9')
curl -sf "http://127.0.0.1:8650/api/printers/eve/jobs/$last/pages" | tee "$work/pages.json"
grep -q '"pages"' "$work/pages.json"
if grep -q '"error"' "$work/pages.json"; then exit 1; fi

# bdf printd: two jobs, one canceled; the server stops right after the
# other is submitted and starts again; the printer has it once. Queues
# of the printer of PostScript and of the printer of the raw port.
server=http://127.0.0.1:8640
cat >"$work/printd.yaml" <<YAML
listen: 127.0.0.1:8640
printers:
  eve: { url: "$url" }
  ps: { url: "$psurl" }
  raw: { url: "socket://127.0.0.1:$rawport" }
queues:
  office: { printers: [eve], defaults: { sides: two-sided-long-edge } }
  postscript: { printers: [ps] }
  old: { printers: [raw] }
fonts: { dirs: ["$root/converter/pptx/testdata/fonts"], no-system: true }
ipp: { listen: ":8633", host: localhost }
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
  if "$bdf" jobs -server "$server" show "$job" | grep '^state: *completed' >/dev/null; then break; fi
  sleep 1
done
"$bdf" jobs -server "$server" -all
"$bdf" jobs -server "$server" show "$job" | grep '^state: *completed' >/dev/null
test "$(received '%PDF-')" -eq $((before + 1))
# the printers of PostScript and of the raw port
psjob=$("$bdf" print -server "$server" -queue postscript -q testdata/form/rows.bdf)
rawjob=$("$bdf" print -server "$server" -queue old -q testdata/form/rows.bdf)
for j in "$psjob" "$rawjob"; do
  for _ in $(seq 120); do
    if "$bdf" jobs -server "$server" show "$j" | grep '^state: *completed' >/dev/null; then break; fi
    sleep 1
  done
  "$bdf" jobs -server "$server" show "$j" | grep '^state: *completed' >/dev/null
done
test "$(received '%!PS-Adobe-3.0' "$psspool")" -eq 3
test "$(received '%!PS-Adobe-3.0' "$rawdir")" -eq 2
draw "$(ls "$rawdir"/3-1.*)"
"$bdf" jobs -server "$server" accounting -csv | tee "$work/accounting.csv"
test "$(wc -l <"$work/accounting.csv")" -eq 5

# bdf printd as an IPP printer: ipptool's checks of IPP/1.1, IPP/2.0 and
# the attributes of IPP Everywhere (ipp-everywhere.test without its
# sample rasters), and its jobs, which reach the printer behind
ippuri="ipp://localhost:8633/ipp/print/office"
tests=/usr/share/cups/ipptool
# (Ubuntu's ipptool has its tests and not their sample files: a PDF and a
# JPEG of the repository stand in)
"$bdf" pdf "${fonts[@]}" testdata/form/rows.bdf "$work/rows.pdf"
cp converter/image/testdata/photo.jpg "$work/photo.jpg"
before=$(received '%PDF-')
ipptool -t -V 2.0 -f "$work/rows.pdf" "$ippuri" test/print/ipp-everywhere-server.test
ipptool -t -f "$work/photo.jpg" "$ippuri" "$tests/print-job.test"
# the jobs of the queue are done
idle() {
  for _ in $(seq 300); do
    if ! "$bdf" jobs -server "$server" -queue office | grep '^ *[0-9]' >/dev/null; then return; fi
    sleep 1
  done
  return 1
}
idle
test "$(received '%PDF-')" -ge $((before + 2))
before=$(received '%PDF-')
# advertised with DNS-SD beside Avahi: bdf printers finds it
"$bdf" printers -timeout 3s | tee "$work/browse.txt"
grep -q '/ipp/print/office' "$work/browse.txt"
# a CUPS queue of the server, as macOS and Linux add it (driverless)
sudo lpadmin -p bdfserver -E -v "$ippuri" -m everywhere
lp -d bdfserver "$work/rows.pdf"
for _ in $(seq 180); do
  if [ "$(received '%PDF-')" -gt "$before" ]; then break; fi
  sleep 1
done
"$bdf" jobs -server "$server" -all -queue office
test "$(received '%PDF-')" -gt "$before"
kill -INT "$printd"
wait "$printd"
printd=
"$bdf" printd -config "$work/printd.yaml" -compact
cat "$work/printd.log"

# labels over 9100 (ZPL), receipts through a raw queue of CUPS (ESC/POS,
# passed on to the same printer of the raw port as it is)
"$bdf" print -printer "socket://127.0.0.1:$rawport?profile=zd421-203" "${fonts[@]}" -media 104x150mm -copies 2 -collate uncollated -darkness 20 testdata/form/rows.bdf
zpl=$(ls -t "$rawdir"/* | head -1)
test "$(head -c 3 "$zpl")" = '^XA'
grep -aqF '^PW831^LL1198' "$zpl"
grep -aqF '^GFA,' "$zpl"
test "$(grep -acF '^PQ2' "$zpl")" -eq 3 # each page of three twice
before=$(received $'\x1b@' "$rawdir")
sudo lpadmin -p bdfraw -E -v "socket://127.0.0.1:$rawport" -m raw
"$bdf" printers -url ipp://localhost:631/printers/bdfraw | grep '^formats:'
"$bdf" print -printer ipp://localhost:631/printers/bdfraw -profile pos-80 "${fonts[@]}" -drawer 1 testdata/form/rows.bdf
for _ in $(seq 60); do
  if [ "$(received $'\x1b@' "$rawdir")" -gt "$before" ]; then break; fi
  sleep 1
done
test "$(received $'\x1b@' "$rawdir")" -gt "$before"
pos=$(ls -t "$rawdir"/* | head -1)
grep -aqF $'\x1dv0' "$pos"
grep -aqF $'\x1dVB' "$pos"
# a form of a template of text: ESC/P in the text mode over 9100, the
# preprinted paper left out
before=$(received $'\x1b@' "$rawdir")
"$bdf" print -printer "socket://127.0.0.1:$rawport?profile=vp-f2000" -skip background testdata/form/text-slip.bdf 2>&1 | grep -F 'ESC/P (text)'
for _ in $(seq 30); do
  if [ "$(received $'\x1b@' "$rawdir")" -gt "$before" ]; then break; fi
  sleep 1
done
escp=$(ls -t "$rawdir"/* | head -1)
grep -aqF $'\x1bC\x21\x1bN\x02' "$escp" # 33 lines, a skip of 2
grep -aqF $'\x1c&' "$escp"                 # Kanji
if grep -aqF $'\x1b*' "$escp"; then exit 1; fi
echo "print: ok"
