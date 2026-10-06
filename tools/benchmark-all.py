#!/usr/bin/env python3
"""Remeasure both workflows on AC power and publish the latest benchmark docs.

Run as your normal user after sudo -v. Results and raw logs stay in /tmp;
only a compact latest summary and the generated documentation go into Git.
"""
import argparse
import hashlib
import importlib.util
import json
import subprocess
import sys
from datetime import datetime
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def run_process_comparison(command):
    # Keep the terminal's sudo credential alive during long Office runs.
    process = subprocess.Popen(command, cwd=ROOT)
    while True:
        try:
            code = process.wait(timeout=60)
            if code:
                raise subprocess.CalledProcessError(code, command)
            return
        except subprocess.TimeoutExpired:
            subprocess.run(['/usr/bin/sudo', '-n', '-v'], capture_output=True, text=True)


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--memory-results', type=Path, help='reuse an AC-powered memory/CPU summary.json')
    p.add_argument('--large-input-dir', type=Path, help='directory containing many-sheets.xlsx and many-slides.pptx')
    p.add_argument('--output', type=Path, help='new raw results directory (default: /tmp/bdf-comparison-TIMESTAMP)')
    p.add_argument('--rounds', type=int, default=3, help='power comparison rounds (default: 3)')
    p.add_argument('--memory-rounds', type=int, default=5)
    p.add_argument('--work-seconds', type=int, default=60)
    a = p.parse_args()
    if a.rounds < 1 or a.memory_rounds < 1 or a.work_seconds < 10:
        p.error('rounds must be positive and work-seconds at least 10')
    state = subprocess.check_output(['pmset', '-g', 'batt'], text=True)
    if 'AC Power' not in state:
        p.error('connect AC power before running the comparison')
    credential = subprocess.run(['/usr/bin/sudo', '-n', '-v'], capture_output=True, text=True)
    if credential.returncode:
        p.error('run sudo -v in this same terminal, then retry')
    output = (a.output or Path('/tmp/bdf-comparison-' + datetime.now().strftime('%Y%m%d-%H%M%S'))).resolve()
    output.mkdir(parents=True, exist_ok=False)
    print(f'Raw results: {output}', flush=True)
    try:
        memory_path = a.memory_results.resolve() if a.memory_results else output/'memory/summary.json'
        if a.memory_results is None:
            command = [sys.executable, str(ROOT/'tools/benchmark-memory.py'), '--rounds', str(a.memory_rounds), '--output', str(output/'memory')]
            for suffix in ('docx', 'pptx', 'xlsx'):
                command += ['--input', str(ROOT/f'examples/sample-files/basic.{suffix}')]
            if a.large_input_dir:
                for name in ('many-sheets.xlsx', 'many-slides.pptx'):
                    source = a.large_input_dir.resolve()/name
                    if not source.is_file():
                        raise ValueError(f'generate the large fixture first: {source}')
                    command += ['--input', str(source)]
            run_process_comparison(command)
        memory = json.loads(memory_path.read_text())
        if not all('AC Power' in memory[k] for k in ('power_source_start', 'power_source_end')):
            raise ValueError('the reused memory measurement must start and end on AC power')
        binary = Path(memory['commands']['bdf'])
        with binary.open('rb') as stream:
            if hashlib.file_digest(stream, 'sha256').hexdigest() != memory['bdf_binary']['sha256']:
                raise ValueError('the measured Go binary has changed; rerun the process comparison')
        subprocess.run([sys.executable, str(ROOT/'tools/benchmark-power.py'), '--bdf', str(binary),
                        '--work-seconds', str(a.work_seconds), '--rounds', str(a.rounds),
                        '--output', str(output/'power')], cwd=ROOT, check=True)
        spec = importlib.util.spec_from_file_location('publisher', ROOT/'tools/update-benchmark-docs.py')
        publisher = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(publisher)
        publisher.publish(memory, json.loads((output/'power/summary.json').read_text()))
        print('Updated docs/benchmarks/latest.json, Japanese/English docs, and charts.', flush=True)
        return 0
    except (OSError, ValueError, KeyError, subprocess.CalledProcessError) as error:
        print(f'Benchmark stopped: {error}\nPartial raw results: {output}', file=sys.stderr)
        return 1


if __name__ == '__main__':
    raise SystemExit(main())
