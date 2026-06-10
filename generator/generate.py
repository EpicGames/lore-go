"""
Copyright Epic Games, Inc. All Rights Reserved.

Generate lorelib wrappers based on Jinja2 template.
"""

import os
import shutil
import sys
import subprocess

from jinja2 import Environment, FileSystemLoader

import common.visitor
from common.generate import generate_templates
from registry import build_augmented
from common.find_lorelib import (
    LIB_DIR,
    MACHINE,
    SYSTEM,
    _wheel_lib_name,
)

# LORE_VERSION_TAG carries the full, un-normalized version (e.g. "v0.8.2-nightly")
# used to build the runtime artifact path; callers that strip LORE_VERSION down to
# a bare "0.8.2" for find_lorelib can pass the original here. Falls back to
# LORE_VERSION.
LORE_VERSION = os.environ.get("LORE_VERSION_TAG") or os.environ.get("LORE_VERSION")
LORE_REVISION = os.environ.get("LORE_REVISION")
LORE_SIBLING_REVISION = os.environ.get("SIBLING_REVISION")
LORE_BRANCH = os.environ.get("LORE_BRANCH")
LORE_SIBLING_BRANCH = os.environ.get("SIBLING_BRANCH")
LORE_NAME = os.environ.get("LORE_NAME")
LORE_RELEASE_BASE_URL = os.environ.get(
    "LORE_RELEASE_BASE_URL",
    "https://github.com/EpicGames/lore/releases/download",
)
# Runtime fetch scheme (see fetch_native_version.ji). Defaults match the
# historical behavior: a raw, short-named direct download.
LORE_ARTIFACT_FORMAT = os.environ.get("LORE_ARTIFACT_FORMAT", "direct")
LORE_ARTIFACT_NAMING = os.environ.get("LORE_ARTIFACT_NAMING", "short")

SCRIPT_DIR = os.path.dirname(__file__)
HEADER_FILE = os.path.join(SCRIPT_DIR, "../lore/include/lore.h")
TEMPLATES_DIR = os.path.join(SCRIPT_DIR, "templates")
SDK_DIR = os.path.join(SCRIPT_DIR, "../lore_go")
TYPES_DIR = os.path.join(SDK_DIR, "types")
NATIVE_DIR = os.path.join(SDK_DIR, "native")
FETCH_NATIVE_DIR = os.path.join(SDK_DIR, "cmd", "fetch-lore-lib")

GENERATE_TARGETS = [
    ("types.ji", TYPES_DIR, "types.go"),
    ("args.ji", TYPES_DIR, "args.go"),
    ("enums.ji", TYPES_DIR, "enums.go"),
    ("events.ji", TYPES_DIR, "events.go"),
    ("functions.ji", NATIVE_DIR, "native.go"),
    ("fluent.ji", SDK_DIR, "lore.go"),
]


def pretty_print_files(generate_targets):
    """Pretty prints the given Go file and updates it in place"""
    for _, directory, file_name in generate_targets:
        content_filename = os.path.join(directory, file_name)
        subprocess.run(["gofmt", "-w", content_filename], check=True)


def copy_lib_to_target_dir():
    """Copies the shared lib binary under target directory"""
    lib_file = _wheel_lib_name()
    src_lib = os.path.join(LIB_DIR, lib_file)

    if SYSTEM == "windows":
        lib_name = "lib/windows/lore.dll"
    elif SYSTEM == "linux":
        if MACHINE in ("arm64", "aarch64"):
            lib_name = "lib/linux_arm64/liblore.so"
        else:
            lib_name = "lib/linux_amd64/liblore.so"
    elif SYSTEM == "darwin":
        lib_name = "lib/darwin/liblore.dylib"
    else:
        sys.exit(f"unsupported platform found: {SYSTEM}")
    dst_lib = os.path.join(SDK_DIR, lib_name)

    os.makedirs(os.path.dirname(dst_lib), exist_ok=True)

    print(f"Copying lore library to {dst_lib}")
    shutil.copy(
        src_lib,
        dst_lib,
    )


generate_templates(
    HEADER_FILE,
    TEMPLATES_DIR,
    GENERATE_TARGETS,
    common.visitor.LoreVisitor,
    build_augmented,
)

print("Applying code formatting", end=" ")
pretty_print_files(GENERATE_TARGETS)
print("done.")

print("Copying the native library", end=" ")
copy_lib_to_target_dir()
print("done.")

# Generate the fetch-lore-lib version file with baked-in Lore version
print("Generating fetch-lore-lib version file", end=" ")
jinja_env = Environment(
    loader=FileSystemLoader(TEMPLATES_DIR), trim_blocks=True, lstrip_blocks=True
)
template = jinja_env.get_template("fetch_native_version.ji")
content = template.render(
    lore_version=LORE_VERSION or "0.0.0",
    lore_revision=LORE_REVISION or "",
    lore_sibling_revision=LORE_SIBLING_REVISION or "",
    lore_branch=LORE_BRANCH or "",
    lore_sibling_branch=LORE_SIBLING_BRANCH or "",
    lore_name=LORE_NAME or "",
    lore_release_base_url=LORE_RELEASE_BASE_URL,
    lore_artifact_format=LORE_ARTIFACT_FORMAT,
    lore_artifact_naming=LORE_ARTIFACT_NAMING,
)
os.makedirs(FETCH_NATIVE_DIR, exist_ok=True)
version_file = os.path.join(FETCH_NATIVE_DIR, "version.go")
with open(version_file, "w", encoding="utf-8") as f:
    f.write(content)
subprocess.run(["gofmt", "-w", version_file], check=True)
print("done.")
