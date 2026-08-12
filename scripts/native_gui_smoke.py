#!/usr/bin/env python3
"""Native Wails GUI smoke tests through GNOME's AT-SPI accessibility tree."""

from __future__ import annotations

import argparse
import os
from pathlib import Path
import shutil
import signal
import subprocess
import tempfile
import time

import gi

gi.require_version("Atspi", "2.0")
from gi.repository import Atspi  # noqa: E402


TIMEOUT_SECONDS = 20
REPOSITORY_ROOT = Path(__file__).resolve().parents[1]


def children(node):
    for index in range(node.get_child_count()):
        yield node.get_child_at_index(index)


def descendants(node):
    pending = [node]
    while pending:
        current = pending.pop()
        yield current
        try:
            pending.extend(reversed(list(children(current))))
        except Exception:
            continue


def wait_for_node(root_supplier, *, name=None, name_contains=None, role=None, timeout=TIMEOUT_SECONDS):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        root = root_supplier()
        if root is not None:
            for node in descendants(root):
                try:
                    node_name = node.get_name() or ""
                    node_role = node.get_role_name()
                    if name is not None and node_name != name:
                        continue
                    if name_contains is not None and name_contains not in node_name:
                        continue
                    if role is not None and node_role != role:
                        continue
                    return node
                except Exception:
                    continue
        time.sleep(0.1)
    description = f"name={name!r}, contains={name_contains!r}, role={role!r}"
    raise AssertionError(f"Timed out waiting for accessible node: {description}")


def find_application(pid):
    desktop = Atspi.get_desktop(0)
    for app in children(desktop):
        try:
            if app.get_name() == "LifeGame" and app.get_process_id() == pid:
                return app
        except Exception:
            continue
    return None


def activate(node):
    count = node.get_n_actions()
    if count < 1:
        raise AssertionError(f"Accessible node has no action: {node.get_name()!r}")
    if not node.do_action(0):
        raise AssertionError(f"Accessible action failed: {node.get_name()!r}")


def launch(binary, home):
    env = os.environ.copy()
    env.update({"HOME": str(home), "NO_AT_BRIDGE": "0"})
    return subprocess.Popen(
        [str(binary)],
        env=env,
        stdout=subprocess.PIPE,
        stderr=subprocess.STDOUT,
        text=True,
    )


def stop(process):
    if process.poll() is not None:
        return
    process.send_signal(signal.SIGINT)
    try:
        process.wait(timeout=5)
    except subprocess.TimeoutExpired:
        process.kill()
        process.wait(timeout=5)


def process_log(process):
    if process.stdout is None:
        return ""
    try:
        return process.stdout.read()
    except Exception:
        return ""


def accessible_snapshot(app):
    if app is None:
        return "<application unavailable>"
    lines = []
    for node in descendants(app):
        try:
            name = node.get_name() or ""
            if name:
                lines.append(f"{node.get_role_name()}: {name}")
        except Exception:
            continue
    return "\n".join(lines)


def assert_user_data_extracted(home):
    data_dir = home / ".lifegame"
    for required_file in (data_dir / "config.yaml", data_dir / "lifegame.db"):
        if not required_file.is_file() or required_file.stat().st_size == 0:
            raise AssertionError(f"Missing initialized user data file: {required_file}")

    expected_root = REPOSITORY_ROOT / "frontend" / "public" / "images"
    extracted_root = data_dir / "images"
    expected_images = {
        path.relative_to(expected_root) for path in expected_root.rglob("*") if path.is_file()
    }
    extracted_images = {
        path.relative_to(extracted_root) for path in extracted_root.rglob("*") if path.is_file()
    }
    if extracted_images != expected_images:
        missing = sorted(str(path) for path in expected_images - extracted_images)
        unexpected = sorted(str(path) for path in extracted_images - expected_images)
        raise AssertionError(
            f"Extracted image tree mismatch; missing={missing[:5]}, unexpected={unexpected[:5]}"
        )


def run_normal_startup(binary, temp_root):
    home = temp_root / "normal-home"
    home.mkdir()
    process = launch(binary, home)
    try:
        app = lambda: find_application(process.pid)
        wait_for_node(app, name="🎮 LifeGame", role="heading")
        wait_for_node(app, name="请输入你的名字", role="entry")
        wait_for_node(app, name="选择女生", role="toggle button")
        wait_for_node(app, name="选择简单难度", role="toggle button")
        wait_for_node(app, name_contains="开始游戏", role="push button")
        assert_user_data_extracted(home)

        activate(wait_for_node(app, name_contains="加载存档", role="push button"))
        wait_for_node(app, name="加载存档", role="dialog")
        activate(wait_for_node(app, name="关闭", role="push button"))
        print("PASS native startup, accessibility tree and dialog interaction")
    except Exception as error:
        snapshot = accessible_snapshot(app())
        stop(process)
        raise AssertionError(f"{error}\n\nAccessibility tree:\n{snapshot}\n\nProcess log:\n{process_log(process)}")
    finally:
        stop(process)


def run_startup_recovery(binary, temp_root):
    blocked_home = temp_root / "blocked-home"
    blocked_home.write_text("not a directory", encoding="utf-8")
    process = launch(binary, blocked_home)
    try:
        app = lambda: find_application(process.pid)
        wait_for_node(app, name="游戏初始化失败", role="heading")
        retry = wait_for_node(app, name="重新尝试", role="push button")

        blocked_home.unlink()
        blocked_home.mkdir()
        activate(retry)
        wait_for_node(app, name="🎮 LifeGame", role="heading")
        print("PASS native startup error page and retry recovery")
    except Exception as error:
        snapshot = accessible_snapshot(app())
        stop(process)
        raise AssertionError(f"{error}\n\nAccessibility tree:\n{snapshot}\n\nProcess log:\n{process_log(process)}")
    finally:
        stop(process)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--binary",
        type=Path,
        default=Path(__file__).resolve().parents[1] / "build" / "bin" / "LifeGame",
    )
    args = parser.parse_args()
    binary = args.binary.resolve()
    if not binary.is_file() or not os.access(binary, os.X_OK):
        raise SystemExit(f"Desktop binary is missing or not executable: {binary}")

    current_desktop = os.environ.get("XDG_CURRENT_DESKTOP", "").lower()
    if "gnome" not in current_desktop:
        raise SystemExit("Native AT-SPI smoke test currently requires a GNOME desktop session")

    original_accessibility = subprocess.check_output(
        ["gsettings", "get", "org.gnome.desktop.interface", "toolkit-accessibility"],
        text=True,
    ).strip()
    temp_root = Path(tempfile.mkdtemp(prefix="lifegame-native-gui-"))

    try:
        if original_accessibility != "true":
            subprocess.run(
                ["gsettings", "set", "org.gnome.desktop.interface", "toolkit-accessibility", "true"],
                check=True,
            )
        Atspi.init()
        run_normal_startup(binary, temp_root)
        run_startup_recovery(binary, temp_root)
    finally:
        if original_accessibility != "true":
            subprocess.run(
                ["gsettings", "set", "org.gnome.desktop.interface", "toolkit-accessibility", original_accessibility],
                check=False,
            )
        shutil.rmtree(temp_root, ignore_errors=True)


if __name__ == "__main__":
    main()
