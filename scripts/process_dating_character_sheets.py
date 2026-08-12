#!/usr/bin/env python3
"""Split generated seven-style dating sheets into runtime WebP assets."""

from __future__ import annotations

import argparse
import subprocess
import tempfile
from pathlib import Path

from PIL import Image, ImageFilter

try:
    import numpy as np
except ImportError:  # Keep the skill helper as a portable fallback.
    np = None


ROOT = Path(__file__).resolve().parents[1]
SHEETS = ROOT / "output/imagegen/dating-fullbody/seven-sheets"
COMPONENTS = ROOT / "output/imagegen/dating-fullbody/seven-components"
PUBLIC = ROOT / "frontend/public/images/datinginfo/dating-partner"
CHROMA_HELPER = Path("/home/parallels/.codex/skills/.system/imagegen/scripts/remove_chroma_key.py")
OUTFITS = ("default", "sleepwear", "romantic", "swimwear", "cosplay", "qipao", "homewear")
OUTFITS_A = ("default", "cosplay", "qipao", "homewear")


def normalize(source: Path, destination: Path) -> None:
    with Image.open(source).convert("RGBA") as image:
        box = image.getbbox()
        if box is None:
            raise ValueError(f"no visible character in {source}")
        crop = image.crop(box)
        scale = min(930 / crop.height, 650 / crop.width)
        crop = crop.resize(
            (max(1, round(crop.width * scale)), max(1, round(crop.height * scale))),
            Image.Resampling.LANCZOS,
        )
        canvas = Image.new("RGBA", (1024, 1024), (0, 0, 0, 0))
        canvas.alpha_composite(crop, ((1024 - crop.width) // 2, 1000 - crop.height))
        destination.parent.mkdir(parents=True, exist_ok=True)
        canvas.save(destination, "WEBP", quality=92, method=4)


def remove_green(source: Path, destination: Path) -> None:
    if np is not None:
        with Image.open(source).convert("RGBA") as image:
            rgba = np.asarray(image, dtype=np.uint8).copy()
        rgb = rgba[..., :3].astype(np.int16)
        border = np.concatenate((rgb[0], rgb[-1], rgb[:, 0], rgb[:, -1]), axis=0)
        key = np.median(border, axis=0).astype(np.int16)
        distance = np.max(np.abs(rgb - key), axis=2).astype(np.float32)
        non_green = np.maximum(rgb[..., 0], rgb[..., 2]).astype(np.float32)
        dominance = rgb[..., 1].astype(np.float32) - non_green
        key_like = (distance <= 32) | (dominance >= 16)

        ratio = np.clip((distance - 12) / (220 - 12), 0, 1)
        soft_alpha = 255 * ratio * ratio * (3 - 2 * ratio)
        denominator = np.maximum(1, float(key[1]) - non_green)
        dominance_alpha = 255 * (1 - np.clip(dominance / denominator, 0, 1))
        alpha = np.where(key_like, np.minimum(soft_alpha, dominance_alpha), 255)
        alpha *= rgba[..., 3].astype(np.float32) / 255
        alpha[(alpha > 0) & (alpha <= 8)] = 0

        spill = key_like & (alpha < 252) & (alpha > 0)
        green_cap = np.maximum(0, non_green - 1).astype(np.uint8)
        rgba[..., 1][spill] = np.minimum(rgba[..., 1][spill], green_cap[spill])
        rgba[..., :3][alpha == 0] = 0
        alpha_image = Image.fromarray(np.rint(alpha).astype(np.uint8), "L").filter(ImageFilter.MinFilter(3))
        rgba[..., 3] = np.asarray(alpha_image)
        Image.fromarray(rgba, "RGBA").save(destination)
        return

    subprocess.run(
        [
            "python3", str(CHROMA_HELPER), "--input", str(source), "--out", str(destination),
            "--auto-key", "border", "--soft-matte", "--transparent-threshold", "12",
            "--opaque-threshold", "220", "--despill", "--edge-contract", "1", "--force",
        ],
        check=True,
    )


def process_sheet(sheet: Path, gender: str, index: int, outfits: tuple[str, ...] = OUTFITS) -> None:
    with Image.open(sheet).convert("RGB") as image, tempfile.TemporaryDirectory() as tmp:
        tmp_dir = Path(tmp)
        edges = [round(image.width * part / len(outfits)) for part in range(len(outfits) + 1)]
        for part, outfit in enumerate(outfits):
            chroma = tmp_dir / f"{outfit}-chroma.png"
            transparent = tmp_dir / f"{outfit}-transparent.png"
            image.crop((edges[part], 0, edges[part + 1], image.height)).save(chroma)
            remove_green(chroma, transparent)
            target_dir = PUBLIC / gender if outfit == "default" else PUBLIC / gender / outfit
            normalize(transparent, target_dir / f"{index:02d}.webp")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--gender", choices=("female", "male"))
    parser.add_argument("--index", type=int)
    args = parser.parse_args()

    completed = 0
    for gender in (args.gender,) if args.gender else ("female", "male"):
        for index in range(1, 51):
            if args.index is not None and index != args.index:
                continue
            component_a = COMPONENTS / gender / f"{index:02d}-a.png"
            component_sleep = COMPONENTS / gender / f"{index:02d}-sleep.png"
            component_romantic = COMPONENTS / gender / f"{index:02d}-romantic.png"
            component_swim = COMPONENTS / gender / f"{index:02d}-swim.png"
            sheet = SHEETS / gender / f"{index:02d}.png"
            component_singles = (component_sleep, component_romantic, component_swim)
            if component_a.exists() and all(item.exists() for item in component_singles):
                process_sheet(component_a, gender, index, OUTFITS_A)
                process_sheet(component_sleep, gender, index, ("sleepwear",))
                process_sheet(component_romantic, gender, index, ("romantic",))
                process_sheet(component_swim, gender, index, ("swimwear",))
            elif sheet.exists():
                process_sheet(sheet, gender, index)
            else:
                continue
            completed += 1
            print(f"processed {gender}/{index:02d}")
    print(f"processed {completed} character sheet(s)")


if __name__ == "__main__":
    main()
