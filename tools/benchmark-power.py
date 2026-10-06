#!/usr/bin/env python3
"""Compare estimated Apple Silicon SoC energy per document.

Run as your normal user after ``sudo -v``. Only powermetrics runs through sudo.
Pass ``--workflow bdf`` to measure only the Go CLI.
Each workflow has an idle phase and a work phase of the same fixed length.
The sampler exits by itself after a finite number of samples, so an interrupted
benchmark cannot leave an indefinitely running privileged sampler behind.

The result covers CPU + GPU + ANE power estimated by powermetrics. It does not
measure total wall power, display, storage, or network equipment. Keep the Mac
on the same power source with a stable screen brightness throughout the run.
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
import time
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SAMPLE_MS = 1000
POWER_LINE = re.compile(
    r"^\s*(Combined Power(?:\s*\([^)]*\))?|CPU Power|GPU Power|ANE Power)"
    r":\s*([\d,.]+)\s*(mW|W)\b",
    re.IGNORECASE | re.MULTILINE,
)


def power_samples(log: str) -> tuple[float, int, str]:
    """Return mean SoC W, sample count, and the selected powermetrics field."""
    found: dict[str, list[float]] = {k: [] for k in ("combined", "cpu", "gpu", "ane")}
    for match in POWER_LINE.finditer(log):
        label, raw, unit = match.groups()
        key = label.split()[0].lower()
        watts = float(raw.replace(",", "")) / (1000 if unit.lower() == "mw" else 1)
        found[key].append(watts)
    if found["combined"]:
        values = found["combined"]
        return statistics.mean(values), len(values), "Combined Power"
    if not found["cpu"]:
        raise ValueError("CPU/Combined Power was not found in powermetrics output")
    counts = [len(values) for values in found.values() if values]
    if max(counts) - min(counts) > 1:
        raise ValueError(f"inconsistent power sample counts: {counts}")
    watts = sum(statistics.mean(values) for values in found.values() if values)
    return watts, min(counts), "+".join(k.upper() for k, v in found.items() if v)


def locate(name: str, supplied: str | None) -> Path:
    if supplied:
        path = Path(supplied).expanduser().resolve()
    else:
        found = shutil.which(name)
        if found:
            path = Path(found).resolve()
        else:
            path = (
                Path.home()
                / ".cache/codex-runtimes/codex-primary-runtime/dependencies/bin/override"
                / name
            )
    if not path.is_file() or not os.access(path, os.X_OK):
        raise ValueError(f"{name} was not found: {path}; pass --{name} explicitly")
    return path


def run_command(args: list[str], cwd: Path = ROOT) -> None:
    result = subprocess.run(args, cwd=cwd, capture_output=True, text=True)
    if result.returncode:
        raise RuntimeError(
            f"command failed ({result.returncode}): {args!r}\n"
            f"stdout:\n{result.stdout[-3000:]}\nstderr:\n{result.stderr[-3000:]}"
        )


def build_bdf(output: Path, supplied: str | None) -> Path:
    if supplied:
        return locate("bdf", supplied)
    binary = output / "bdf"
    print("Building the release bdf CLI (outside the measurement)...", flush=True)
    run_command(
        ["go", "build", "-trimpath", "-ldflags=-s -w", "-o", str(binary), "./cmd/bdf"]
    )
    return binary


def sha256(path: Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def describe(args: list[str]) -> str:
    try:
        result = subprocess.run(args, capture_output=True, text=True, timeout=15)
        return (result.stdout + result.stderr).strip()[:500]
    except (OSError, subprocess.TimeoutExpired) as error:
        return str(error)


def convert_one(workflow: str, source: Path, dest: Path, commands: dict[str, Path]) -> None:
    dest.mkdir(parents=True)
    if workflow == "libreoffice":
        # A private profile prevents an existing GUI/daemon from accepting the job.
        profile = dest / "lo-profile"
        pdf = dest / (source.stem + ".pdf")
        run_command(
            [
                str(commands["soffice"]),
                f"-env:UserInstallation={profile.as_uri()}",
                "--headless",
                "--convert-to",
                "pdf",
                "--outdir",
                str(dest),
                str(source),
            ]
        )
        if not pdf.is_file() or not pdf.stat().st_size:
            raise RuntimeError(f"LibreOffice did not create {pdf}")
        run_command(
            [
                str(commands["pdftoppm"]),
                "-f",
                "1",
                "-singlefile",
                "-scale-to",
                "256",
                "-png",
                str(pdf),
                str(dest / "thumb"),
            ]
        )
        thumb = dest / "thumb.png"
    else:
        thumb = dest / "thumb.png"
        document = dest / "document.bdf"
        run_command(
            [
                str(commands["bdf"]),
                "generate",
                "-q",
                "-thumbnail",
                str(thumb),
                "-thumbnail-mode",
                "fit",
                str(source),
                str(document),
            ]
        )
        if not document.is_file() or not document.stat().st_size:
            raise RuntimeError(f"bdf did not create {document}")
    if not thumb.is_file() or not thumb.stat().st_size:
        raise RuntimeError(f"thumbnail was not created: {thumb}")


def measure_phase(
    output: Path, name: str, seconds: int, action=None
) -> dict[str, object]:
    # Refresh the sudo timestamp before every finite phase. A full run takes
    # longer than macOS's usual sudo timeout, but no password prompt occurs
    # during a power sample.
    credential = subprocess.run(
        ["/usr/bin/sudo", "-n", "-v"], capture_output=True, text=True
    )
    if credential.returncode:
        raise RuntimeError(
            "sudo authentication is unavailable or expired. Run sudo -v "
            "in this terminal, then start the benchmark again."
        )
    log_path = output / f"{name}.powermetrics.txt"
    error_path = output / f"{name}.powermetrics.stderr.txt"
    command = [
        "/usr/bin/sudo", "-n", "/usr/bin/powermetrics",
        "-s", "cpu_power,gpu_power,ane_power",
        "-i", str(SAMPLE_MS), "-n", str(seconds), "-f", "text", "-b", "1", "-a", "0",
    ]
    print(f"  {name}: {seconds} s", flush=True)
    with log_path.open("wb") as log, error_path.open("wb") as errors:
        started = time.monotonic()
        sampler = subprocess.Popen(command, stdout=log, stderr=errors)
        action_started = action_ended = None
        try:
            time.sleep(2)
            if sampler.poll() is not None:
                raise RuntimeError(
                    f"powermetrics exited early; read {error_path}. "
                    "Run sudo -v in this terminal, then retry."
                )
            if action is not None:
                action_started = time.monotonic()
                action()
                action_ended = time.monotonic()
                if sampler.poll() is not None:
                    raise RuntimeError(
                        f"The workload exceeded the {seconds}-second sampling window. "
                        "Increase --libreoffice-seconds or --bdf-seconds."
                    )
            sampler.wait(timeout=seconds + 20)
        except BaseException:
            # powermetrics also has a finite -n count if sudo cannot forward SIGINT.
            if sampler.poll() is None:
                sampler.send_signal(2)
            raise
        ended = time.monotonic()
    if sampler.returncode:
        raise RuntimeError(f"powermetrics failed ({sampler.returncode}); read {error_path}")
    average, count, field = power_samples(log_path.read_text(errors="replace"))
    if count < max(2, seconds // 2):
        raise RuntimeError(f"Only {count} power samples in {log_path}; expected about {seconds}")
    result: dict[str, object] = {
        "name": name,
        "seconds_requested": seconds,
        "seconds_actual": round(ended - started, 3),
        "samples": count,
        "power_field": field,
        "mean_soc_w": round(average, 6),
        "log": str(log_path),
        "stderr": str(error_path),
    }
    if action_started is not None and action_ended is not None:
        result["work_seconds"] = round(action_ended - action_started, 3)
    return result


def save(output: Path, data: dict[str, object]) -> None:
    path = output / "summary.json"
    temp = output / "summary.json.tmp"
    temp.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n")
    temp.replace(path)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", type=Path, help="new result directory (default: /tmp/bdf-power-TIMESTAMP)")
    parser.add_argument("--bdf", help="prebuilt bdf CLI; otherwise build a release binary")
    parser.add_argument("--soffice", help="LibreOffice/OpenOffice soffice executable")
    parser.add_argument("--pdftoppm", help="Poppler pdftoppm executable")
    parser.add_argument("--input", action="append", type=Path, help="repeat to replace the three default samples")
    parser.add_argument("--workflow", action="append", choices=("bdf", "libreoffice"),
                        help="repeat to select workflows (default: both); bdf measures Go only")
    parser.add_argument("--jobs", type=int, default=90, help="documents in each work phase (default: 90)")
    parser.add_argument("--rounds", type=int, default=2, help="paired rounds (default: 2)")
    parser.add_argument("--libreoffice-seconds", type=int, default=120, help="samples per LibreOffice phase (default: 120)")
    parser.add_argument("--bdf-seconds", type=int, default=60, help="samples per bdf phase (default: 60)")
    args = parser.parse_args()
    workflows = list(dict.fromkeys(args.workflow or ("libreoffice", "bdf")))
    phase_seconds = {"libreoffice": args.libreoffice_seconds, "bdf": args.bdf_seconds}
    if args.jobs < 1 or args.rounds < 1 or any(phase_seconds[w] < 10 for w in workflows):
        parser.error("jobs and rounds must be positive; each sampling phase must last at least 10 seconds")
    if platform.system() != "Darwin":
        parser.error("this benchmark requires macOS powermetrics")
    if os.geteuid() == 0:
        parser.error("run this script as your normal user; it calls sudo only for powermetrics")

    sources = [path.expanduser().resolve() for path in (args.input or [
        ROOT / "examples/sample-files/basic.docx",
        ROOT / "examples/sample-files/basic.pptx",
        ROOT / "examples/sample-files/basic.xlsx",
    ])]
    for source in sources:
        if not source.is_file():
            parser.error(f"input does not exist: {source}")
    commands = {}
    if "libreoffice" in workflows:
        commands["soffice"] = locate("soffice", args.soffice)
        commands["pdftoppm"] = locate("pdftoppm", args.pdftoppm)
    stamp = datetime.now().strftime("%Y%m%d-%H%M%S")
    output = (args.output or Path(f"/tmp/bdf-power-{stamp}")).expanduser().resolve()
    output.mkdir(parents=True, exist_ok=False)
    try:
        if "bdf" in workflows:
            commands["bdf"] = build_bdf(output, args.bdf)
        data: dict[str, object] = {
            "created_utc": datetime.now(timezone.utc).isoformat(),
            "host": platform.platform(),
            "scope": "estimated CPU+GPU+ANE SoC energy, not whole-device energy",
            "jobs_per_phase": args.jobs,
            "rounds": args.rounds,
            "workflows": workflows,
            "phase_seconds": {w: phase_seconds[w] for w in workflows},
            "sample_interval_ms": SAMPLE_MS,
            "commands": {name: str(path) for name, path in commands.items()},
            "versions": {},
            "power_source_start": describe(["pmset", "-g", "batt"]),
            "inputs": [{"path": str(p), "bytes": p.stat().st_size, "sha256": sha256(p)} for p in sources],
            "phases": [],
            "pairs": [],
        }
        if "bdf" in commands:
            data["bdf_binary"] = {
                "bytes": commands["bdf"].stat().st_size,
                "sha256": sha256(commands["bdf"]),
            }
            data["versions"]["go"] = describe(["go", "version"])
        if "soffice" in commands:
            data["versions"]["soffice"] = describe([str(commands["soffice"]), "--version"])
            data["versions"]["pdftoppm"] = describe([str(commands["pdftoppm"]), "-v"])
        save(output, data)
        print(f"Results: {output}", flush=True)
        print("Checking powermetrics output and sudo access (3 s)...", flush=True)
        preflight = measure_phase(output, "preflight", 3)
        data["preflight"] = preflight
        save(output, data)
        print(f"  {preflight['power_field']}: {preflight['mean_soc_w']} W", flush=True)

        print("Warming up file caches outside the measurements...", flush=True)
        with tempfile.TemporaryDirectory(prefix="bdf-power-warmup-", dir=output) as warmup:
            for workflow in workflows:
                for i, source in enumerate(sources):
                    convert_one(workflow, source, Path(warmup) / workflow / str(i), commands)
        time.sleep(5)

        for round_index in range(args.rounds):
            order = workflows if round_index % 2 == 0 else list(reversed(workflows))
            for workflow in order:
                seconds = phase_seconds[workflow]
                label = f"round{round_index + 1}-{workflow}"
                idle = measure_phase(output, label + "-idle", seconds)
                data["phases"].append(idle)
                save(output, data)
                with tempfile.TemporaryDirectory(prefix=label + "-", dir=output) as workdir:
                    directory = Path(workdir)

                    def work() -> None:
                        for i in range(args.jobs):
                            convert_one(workflow, sources[i % len(sources)], directory / f"job-{i:04d}", commands)

                    active = measure_phase(output, label + "-work", seconds, work)
                    data["phases"].append(active)
                    watts = float(active["mean_soc_w"]) - float(idle["mean_soc_w"])
                    joules_per_job = watts * float(active["seconds_actual"]) / args.jobs
                    pair = {
                        "round": round_index + 1,
                        "workflow": workflow,
                        "estimated_incremental_soc_w": round(watts, 6),
                        "estimated_soc_j_per_document": round(joules_per_job, 6),
                        "idle": idle["name"],
                        "work": active["name"],
                    }
                    data["pairs"].append(pair)
                    save(output, data)
                    print(
                        f"  {label}: {joules_per_job:.3f} estimated SoC J/document "
                        f"({active['work_seconds']} s work within {active['seconds_actual']} s sampling)",
                        flush=True,
                    )
                time.sleep(5)  # keep output cleanup out of the next idle sample

        medians = {
            workflow: statistics.median(
                float(pair["estimated_soc_j_per_document"])
                for pair in data["pairs"] if pair["workflow"] == workflow
            )
            for workflow in workflows
        }
        data["median_estimated_soc_j_per_document"] = medians
        data["power_source_end"] = describe(["pmset", "-g", "batt"])
        if medians.get("libreoffice", 0) > 0 and medians.get("bdf", -1) >= 0:
            data["estimated_soc_energy_reduction_percent"] = round(
                100 * (1 - medians["bdf"] / medians["libreoffice"]), 2
            )
        save(output, data)
        print(json.dumps({
            "median_estimated_soc_j_per_document": medians,
            "estimated_soc_energy_reduction_percent": data.get("estimated_soc_energy_reduction_percent"),
            "summary": str(output / "summary.json"),
        }, indent=2), flush=True)
        return 0
    except (OSError, RuntimeError, ValueError, subprocess.TimeoutExpired) as error:
        print(f"Benchmark stopped: {error}", file=sys.stderr)
        print(f"Any raw logs and partial summary remain in {output}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    raise SystemExit(main())
