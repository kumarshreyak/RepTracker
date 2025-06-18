#!/usr/bin/env python3
"""
Scrape StrengthLog’s exercise directory and save the data in JSON.

Install first:
    pip install requests beautifulsoup4 tqdm
Run:
    python scrape_strengthlog.py
"""
import json, re, sys, time, hashlib
from datetime import datetime, timezone
from pathlib import Path
from typing import List, Dict, Tuple

import requests
from bs4 import BeautifulSoup
from tqdm import tqdm


BASE_URL  = "https://www.strengthlog.com"
DIR_URL   = f"{BASE_URL}/exercise-directory/"

OUTPUT    = Path("strengthlog_exercises.json")
TIMESTAMP = datetime.now(timezone.utc).isoformat(timespec="seconds")


def slugify(text: str) -> str:
    return re.sub(r"[^a-z0-9]+", "-", text.lower()).strip("-")


def guess_equipment(name: str) -> str:
    # Very simple heuristics – tweak as needed.
    for word in ("Barbell", "Dumbbell", "Kettlebell", "Machine",
                 "Smith Machine", "Cable", "Resistance Band", "Band",
                 "Bodyweight", "Rowing Machine", "Stationary Bike",
                 "Trap Bar", "Landmine"):
        if word.lower() in name.lower():
            return word
    return "Bodyweight"


def guess_difficulty(name: str) -> str:
    # Cheap heuristic: compound barbell lifts → Advanced, otherwise Intermediate
    advanced = ("Deadlift", "Squat", "Clean", "Snatch", "Jerk", "Bench Press")
    return "Advanced" if any(w.lower() in name.lower() for w in advanced) else "Intermediate"


def parse_instructions(soup: BeautifulSoup) -> List[str]:
    """
    Pull the first ordered/numbered list OR the first 5 bullet points that
    look like how-to steps.  Falls back to empty list if nothing obvious.
    """
    # 1. Try a <ol> list
    ol = soup.find("ol")
    if ol:
        return [li.get_text(" ", strip=True) for li in ol.find_all("li")]

    # 2. Try <ul> bullet list with 5-10 items
    ul = soup.find("ul")
    if ul:
        bullets = [li.get_text(" ", strip=True) for li in ul.find_all("li")]
        if 2 <= len(bullets) <= 15:
            return bullets

    return []


def scrape_exercise_page(url: str) -> Tuple[str, List[str]]:
    """Return description paragraph + how-to steps (if any)."""
    res  = requests.get(url, timeout=30)
    soup = BeautifulSoup(res.text, "html.parser")

    # First non-empty paragraph below the H1
    hdr = soup.find("h1")
    para = hdr.find_next("p") if hdr else soup.find("p")
    description = para.get_text(" ", strip=True) if para else ""

    instructions = parse_instructions(soup)

    return description, instructions


def main() -> None:
    print("Fetching directory …", file=sys.stderr)
    res   = requests.get(DIR_URL, timeout=30)
    soup  = BeautifulSoup(res.text, "html.parser")

    data: List[Dict[str, object]] = []
    next_id = 1

    # The directory is structured as h3 (“### Chest Exercises”) headings
    # followed by an ordered list.  Traverse every h3 that ends with “Exercises”.
    for h3 in soup.find_all("h3"):
        title = h3.get_text(strip=True)
        if not title.endswith("Exercises"):
            continue

        muscle_group = title.replace(" Exercises", "")
        ol = h3.find_next("ol")
        if not ol:
            continue

        for li in ol.find_all("li"):
            # Each list item is: “1. Exercise Name”
            a = li.find("a")
            if not a:
                continue
            name = a.get_text(strip=True)
            url  = a["href"]

            desc, steps = scrape_exercise_page(url)

            exercise = {
                "id": str(next_id),
                "name": name,
                "description": desc or f"{name} exercise.",
                "muscleGroup": muscle_group,
                "equipment": guess_equipment(name),
                "difficulty": guess_difficulty(name),
                "instructions": steps,
                "createdAt": TIMESTAMP,
                "updatedAt": TIMESTAMP,
            }
            data.append(exercise)
            next_id += 1

    OUTPUT.write_text(json.dumps(data, indent=2, ensure_ascii=False))
    print(f"Wrote {len(data)} exercises to {OUTPUT}", file=sys.stderr)


if __name__ == "__main__":
    main()
