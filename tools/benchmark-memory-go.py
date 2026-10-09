#!/usr/bin/env python3
"""Measure Go conversion + 256px thumbnail process peaks on macOS."""

import argparse
import hashlib
import json
import platform
import re
import statistics
import subprocess
import tempfile
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
PATTERNS = {
    "max_rss_bytes": r"^\s*(\d+)\s+maximum resident set size\s*$",
    "peak_footprint_bytes": r"^\s*(\d+)\s+peak memory footprint\s*$",
    "elapsed_seconds": r"^\s*([\d.]+) real\s+[\d.]+ user\s+[\d.]+ sys\s*$",
    "user_cpu_seconds": r"^\s*[\d.]+ real\s+([\d.]+) user\s+[\d.]+ sys\s*$",
    "system_cpu_seconds": r"^\s*[\d.]+ real\s+[\d.]+ user\s+([\d.]+) sys\s*$",
}


def sha256(path):
    with path.open("rb") as stream:
        return hashlib.file_digest(stream, "sha256").hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--bdf", required=True, type=Path)
    parser.add_argument("--compare-bdf", type=Path, help="also measure a pre-change binary, alternating order")
    parser.add_argument("--compare-revision", help="source revision of the comparison binary")
    parser.add_argument("--baseline", type=Path, help="optional historical LibreOffice comparison JSON")
    parser.add_argument("--input", required=True, action="append", type=Path)
    parser.add_argument("--rounds", type=int, default=5)
    parser.add_argument("--revision", help="source revision of a prebuilt binary")
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    if platform.system() != "Darwin" or args.rounds < 1:
        parser.error("macOS and a positive round count are required")
    binary = args.bdf.resolve()
    baseline = json.loads(args.baseline.read_text()) if args.baseline else None
    data = {
        "created_utc": datetime.now(timezone.utc).isoformat(),
        "git_revision": args.revision or subprocess.check_output(
            ["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip(),
        "host": platform.platform(),
        "go": subprocess.check_output(["go", "version"], text=True).strip(),
        "method": "/usr/bin/time -l; fresh process per conversion; document plus 256px fit thumbnail; one excluded warmup per input; medians",
        "rounds": args.rounds,
        "build_flags": ["-trimpath", "-ldflags=-s -w"],
        "binary": {"path": str(binary), "bytes": binary.stat().st_size, "sha256": sha256(binary)},
        "power_source_start": subprocess.check_output(["pmset", "-g", "batt"], text=True).strip(),
        "baseline": {"date": baseline["date"], "git_revision": baseline["git_revision"],
                     "path": str(args.baseline.resolve()), "sha256": sha256(args.baseline)} if baseline else None,
        "inputs": [], "runs": [], "results": {},
        "notes": ["LibreOffice and Poppler are historical results, not remeasured.",
                  "Large inputs were regenerated with the original recipe and package versions; ZIP hashes can differ.",
                  "Historical many-slides bdf_only results do not record the command or elapsed time; that comparison is a reference only."],
    }
    binaries = {"current": binary}
    if args.compare_bdf:
        before = args.compare_bdf.resolve()
        binaries["before"] = before
        data["comparison_binary"] = {"path": str(before), "bytes": before.stat().st_size,
                                     "sha256": sha256(before), "git_revision": args.compare_revision}
        data["method"] += "; before/current order alternates each round; PNG bytes compared"
    if not baseline:
        data["notes"] = []
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.TemporaryDirectory(prefix="bdf-go-process-benchmark-") as temp:
        for source in map(Path.resolve, args.input):
            name = source.name
            previous, old, office, old_time, office_time = {}, {}, None, None, None
            if baseline:
                previous = next((f for f in baseline["small_inputs"]["files"] if f["name"] == name), None)
                if previous:
                    historic = baseline["small_inputs"]["results"][name]
                    old = historic["median_pipeline_peak_bytes"]["bdf"]
                    office = historic["median_pipeline_peak_bytes"]["libreoffice"]
                    old_time = historic["median_pipeline_elapsed_seconds"]["bdf"]
                    office_time = historic["median_pipeline_elapsed_seconds"]["libreoffice"]
                else:
                    previous = baseline["large_inputs"]["files"][name]
                    historic = baseline["large_inputs"]["xlsx_pipeline"] if name.endswith("xlsx") else None
                    old = historic["median_bytes"]["bdf"] if historic else baseline["large_inputs"]["bdf_only"]["median_bytes"][name]
                    office = historic["median_bytes"]["libreoffice"] if historic else None
                    old_time = historic["median_pipeline_elapsed_seconds"]["bdf"] if historic else None
                    office_time = historic["median_pipeline_elapsed_seconds"]["libreoffice"] if historic else None
            digest = sha256(source)
            if previous and name.startswith("basic.") and digest != previous["sha256"]:
                raise RuntimeError(f"historical input hash mismatch: {source}")
            data["inputs"].append({"path": str(source), "bytes": source.stat().st_size,
                                   "sha256": digest, "historical_sha256": previous.get("sha256"),
                                   "matches_historical_hash": digest == previous["sha256"] if previous else None})
            runs = []
            for index in range(args.rounds + 1):
                order = list(binaries) if index % 2 == 0 else list(reversed(binaries))
                thumbnails = []
                for label in order:
                    dest = Path(temp) / f"{name}-{index}-{label}"
                    dest.mkdir()
                    command = [str(binaries[label]), "generate", "-q", "-thumbnail", str(dest / "thumb.png"),
                               "-thumbnail-mode", "fit", str(source), str(dest / "document.bdf")]
                    proc = subprocess.run(["/usr/bin/time", "-l", *command], cwd=ROOT,
                                          capture_output=True, text=True, check=True)
                    values = {}
                    for metric, pattern in PATTERNS.items():
                        match = re.search(pattern, proc.stderr, re.M)
                        if not match:
                            raise RuntimeError(f"{metric} missing: {proc.stderr}")
                        values[metric] = float(match[1]) if metric.endswith("_seconds") else int(match[1])
                    values["cpu_seconds"] = round(values["user_cpu_seconds"] + values["system_cpu_seconds"], 4)
                    for output in [dest / "thumb.png", dest / "document.bdf"]:
                        if not output.is_file() or not output.stat().st_size:
                            raise RuntimeError(f"missing output: {output}")
                    values["output_bytes"] = (dest / "document.bdf").stat().st_size
                    thumbnails.append((dest / "thumb.png").read_bytes())
                    if index:
                        run = {"input": name, "round": index, "variant": label, "command": command, **values,
                               "thumbnail_sha256": sha256(dest / "thumb.png"), "raw_time_output": proc.stderr}
                        runs.append(run)
                        data["runs"].append(run)
                if any(png != thumbnails[0] for png in thumbnails):
                    raise RuntimeError(f"thumbnail changed for {name}, round {index}")
                print(f"{name}: {'warmup' if index == 0 else f'round {index}'}", flush=True)
            result = {"historical_bdf": {**old, "elapsed_seconds": old_time} if old else None,
                      "historical_libreoffice": {**office, "elapsed_seconds": office_time} if office else None}
            for label in binaries:
                result[label] = {metric: statistics.median(r[metric] for r in runs if r["variant"] == label)
                                 for metric in [*PATTERNS, "cpu_seconds", "output_bytes"]}
            result["thumbnails_identical"] = True if args.compare_bdf else None
            data["results"][name] = result
            args.output.write_text(json.dumps(data, indent=2) + "\n")
    data["power_source_end"] = subprocess.check_output(["pmset", "-g", "batt"], text=True).strip()
    args.output.write_text(json.dumps(data, indent=2) + "\n")
    print(json.dumps(data["results"], indent=2))


if __name__ == "__main__":
    main()
