#!/usr/bin/env python3
"""Measure peak memory for Office PDF+thumbnail and bdf+thumbnail workflows.

macOS only. Uses /usr/bin/time -l; no administrator access is needed. Each
LibreOffice conversion starts with a private profile. LibreOffice and Poppler
run sequentially, so the workflow peak is the larger of their two peaks.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import platform
import re
import shutil
import statistics
import subprocess
import sys
import tempfile
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
METRICS = {
    "max_rss_bytes": re.compile(r"^\s*(\d+)\s+maximum resident set size\s*$", re.M),
    "peak_footprint_bytes": re.compile(r"^\s*(\d+)\s+peak memory footprint\s*$", re.M),
}
ELAPSED = re.compile(r"^\s*([\d.]+) real\s+([\d.]+) user\s+([\d.]+) sys\s*$", re.M)


def locate(name: str, supplied: str | None) -> Path:
    if supplied:
        path = Path(supplied).expanduser().resolve()
    else:
        found = shutil.which(name)
        path = Path(found).resolve() if found else (
            Path.home()
            / ".cache/codex-runtimes/codex-primary-runtime/dependencies/bin/override"
            / name
        )
    if not path.is_file() or not os.access(path, os.X_OK):
        raise ValueError(f"{name} was not found: {path}; pass --{name} explicitly")
    return path


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def version(command: list[str]) -> str:
    result = subprocess.run(command, capture_output=True, text=True, timeout=15)
    return (result.stdout + result.stderr).strip()


def timed(command: list[str], log: Path) -> dict[str, int | float]:
    result = subprocess.run(
        ["/usr/bin/time", "-l", *command], cwd=ROOT, capture_output=True, text=True
    )
    log.write_text(
        f"command: {command!r}\nexit: {result.returncode}\n"
        f"stdout:\n{result.stdout}\nstderr:\n{result.stderr}"
    )
    if result.returncode:
        raise RuntimeError(f"command failed ({result.returncode}); see {log}")
    values = {}
    for name, pattern in METRICS.items():
        match = pattern.search(result.stderr)
        if not match:
            raise RuntimeError(f"{name} missing from {log}")
        values[name] = int(match.group(1))
    elapsed = ELAPSED.search(result.stderr)
    if not elapsed:
        raise RuntimeError(f"elapsed time missing from {log}")
    values["elapsed_seconds"] = float(elapsed.group(1))
    values["cpu_seconds"] = round(float(elapsed.group(2)) + float(elapsed.group(3)), 4)
    return values


def convert(workflow: str, source: Path, dest: Path, commands: dict[str, Path]) -> dict:
    dest.mkdir(parents=True, exist_ok=True)
    if workflow == "libreoffice":
        pdf = dest / f"{source.stem}.pdf"
        stages = {
            "soffice": timed([
                str(commands["soffice"]),
                f"-env:UserInstallation={(dest / 'lo-profile').as_uri()}",
                "--headless", "--convert-to", "pdf", "--outdir", str(dest), str(source),
            ], dest / "soffice.time.txt")
        }
        if not pdf.is_file() or not pdf.stat().st_size:
            raise RuntimeError(f"PDF missing: {pdf}")
        stages["pdftoppm"] = timed([
            str(commands["pdftoppm"]), "-f", "1", "-singlefile", "-scale-to", "256",
            "-png", str(pdf), str(dest / "thumb"),
        ], dest / "pdftoppm.time.txt")
    else:
        document = dest / "document.bdf"
        stages = {"bdf": timed([
            str(commands["bdf"]), "generate", "-q", "-thumbnail",
            str(dest / "thumb.png"), "-thumbnail-mode", "fit",
            str(source), str(document),
        ], dest / "bdf.time.txt")}
        if not document.is_file() or not document.stat().st_size:
            raise RuntimeError(f"bdf missing: {document}")
    if not (dest / "thumb.png").is_file():
        raise RuntimeError(f"thumbnail missing: {dest / 'thumb.png'}")
    return {
        "stages": stages,
        "pipeline_elapsed_seconds": round(sum(stage["elapsed_seconds"] for stage in stages.values()), 2),
        "pipeline_cpu_seconds": round(sum(stage["cpu_seconds"] for stage in stages.values()), 4),
        "pipeline_peak": {
            metric: max(stage[metric] for stage in stages.values())
            for metric in METRICS
        },
    }


def save(path: Path, data: dict) -> None:
    temporary = path.with_suffix(".tmp")
    temporary.write_text(json.dumps(data, indent=2) + "\n")
    temporary.replace(path)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, help="new results directory")
    parser.add_argument("--bdf", help="prebuilt bdf CLI; otherwise build a release binary")
    parser.add_argument("--soffice", help="LibreOffice/OpenOffice soffice executable")
    parser.add_argument("--pdftoppm", help="Poppler pdftoppm executable")
    parser.add_argument("--input", action="append", type=Path, help="repeat for custom samples")
    parser.add_argument("--rounds", type=int, default=5, help="measured runs after warmup (default: 5)")
    args = parser.parse_args()
    if platform.system() != "Darwin":
        parser.error("this benchmark requires macOS /usr/bin/time -l")
    if args.rounds < 1:
        parser.error("--rounds must be positive")
    sources = [path.expanduser().resolve() for path in (args.input or [
        ROOT / "examples/sample-files/basic.docx",
        ROOT / "examples/sample-files/basic.pptx",
        ROOT / "examples/sample-files/basic.xlsx",
    ])]
    for source in sources:
        if not source.is_file():
            parser.error(f"input does not exist: {source}")
    commands = {
        "soffice": locate("soffice", args.soffice),
        "pdftoppm": locate("pdftoppm", args.pdftoppm),
    }
    stamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    output = (args.output or Path(f"/tmp/bdf-memory-{stamp}")).expanduser().resolve()
    output.mkdir(parents=True, exist_ok=False)
    try:
        commands["bdf"] = locate("bdf", args.bdf) if args.bdf else output / "bdf"
        if not args.bdf:
            print("Building the release bdf CLI (outside the measurement)...", flush=True)
            subprocess.run(
                ["go", "build", "-trimpath", "-ldflags=-s -w", "-o",
                 str(commands["bdf"]), "./cmd/bdf"], cwd=ROOT, check=True
            )
        data = {
            "created_utc": datetime.now(timezone.utc).isoformat(),
            "host": platform.platform(),
            "method": "/usr/bin/time -l; per-process peak; sequential pipeline peak is max of stages",
            "rounds": args.rounds,
            "power_source_start": subprocess.check_output(["pmset", "-g", "batt"], text=True).strip(),
            "versions": {name: version([str(path), "-v" if name == "pdftoppm" else "--version"])
                         for name, path in commands.items() if name != "bdf"},
            "commands": {name: str(path) for name, path in commands.items()},
            "bdf_binary": {"bytes": commands["bdf"].stat().st_size,
                           "sha256": sha256(commands["bdf"])},
            "inputs": [{"path": str(p), "bytes": p.stat().st_size,
                        "sha256": sha256(p)} for p in sources],
            "runs": [],
        }
        data["versions"]["go"] = version(["go", "version"])
        summary_path = output / "summary.json"
        save(summary_path, data)
        for source in sources:
            for round_index in range(args.rounds + 1):
                warmup = round_index == 0
                order = ("libreoffice", "bdf") if round_index % 2 == 0 else ("bdf", "libreoffice")
                for workflow in order:
                    label = f"{source.name}-{round_index:02d}-{workflow}"
                    with tempfile.TemporaryDirectory(prefix=label + "-", dir=output) as temp:
                        result = convert(workflow, source, Path(temp), commands)
                        # Keep command and measurement logs, but discard all converted outputs.
                        for log in Path(temp).glob("*.time.txt"):
                            shutil.copy2(log, output / f"{label}.{log.name}")
                    if not warmup:
                        data["runs"].append({
                            "input": source.name, "round": round_index,
                            "workflow": workflow, **result,
                        })
                        save(summary_path, data)
                print(f"{source.name}: {'warmup' if warmup else f'round {round_index}'}", flush=True)
        data["median_pipeline_peak_bytes"] = {
            source.name: {
                workflow: {
                    metric: statistics.median(
                        run["pipeline_peak"][metric] for run in data["runs"]
                        if run["input"] == source.name and run["workflow"] == workflow
                    )
                    for metric in METRICS
                }
                for workflow in ("libreoffice", "bdf")
            }
            for source in sources
        }
        data["median_pipeline_elapsed_seconds"] = {
            source.name: {
                workflow: statistics.median(
                    run["pipeline_elapsed_seconds"] for run in data["runs"]
                    if run["input"] == source.name and run["workflow"] == workflow
                )
                for workflow in ("libreoffice", "bdf")
            }
            for source in sources
        }
        data["median_pipeline_cpu_seconds"] = {
            source.name: {
                workflow: statistics.median(
                    run["pipeline_cpu_seconds"] for run in data["runs"]
                    if run["input"] == source.name and run["workflow"] == workflow
                ) for workflow in ("libreoffice", "bdf")
            } for source in sources
        }
        data["power_source_end"] = subprocess.check_output(["pmset", "-g", "batt"], text=True).strip()
        save(summary_path, data)
        print(json.dumps({"median_pipeline_peak_bytes": data["median_pipeline_peak_bytes"],
                          "median_pipeline_elapsed_seconds": data["median_pipeline_elapsed_seconds"],
                          "summary": str(summary_path)}, indent=2), flush=True)
        return 0
    except (OSError, RuntimeError, subprocess.CalledProcessError) as error:
        print(f"Benchmark stopped: {error}\nPartial logs: {output}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
