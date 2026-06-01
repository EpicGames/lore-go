"""
Copyright Epic Games, Inc. All Rights Reserved.

Type registries used by the Jinja templates. The seed maps below are the
hand-curated mappings for primitive, FFI, and custom-builder types. The
seed lists track types that are either hand-written (skipped by the main
type-gen loop) or filtered out from the generated function surface.

`build_augmented` extends the seeds with auto-detected entries for every
`*_array_t` typedef found in `lore.h` (split into `array_types` for
regular arrays and `event_array_types` for event-data arrays), so adding
a new array type to the header requires zero edits here or in the Jinja
templates.
"""

from common import util


SEED_IGNORED_FUNCTIONS = [
    "lore_event_type",
    "lore_shutdown",
    "lore_set_allocator",
    "lore_version",
    "lore_user_directory",
    "lore_log_configure",
]


SEED_CUSTOM_TYPES = [
    "lore_string_t",
    "lore_hash_t",
    "lore_context_t",
    "lore_partition_t",
    "lore_address_t",
    "lore_branch_id_t",
    "lore_repository_id_t",
    "lore_instance_id_t",
    "lore_string_array_t",
    "lore_uint8_array_t",
    "lore_event_t",
    "lore_metadata_t",
    "lore_event_callback_config_t",
    "lore_binary_t",
    "lore_bytes_t",
]


SEED_C_TO_GO_FFI_TYPE_MAP = {
    "lore_metadata_t": "LoreMetadataFFI",
    "lore_binary_t": "LoreBinaryFFI",
    "lore_bytes_t": "LoreBytesFFI",
}


SEED_C_TO_GO_TYPE_MAP = {
    "uintptr_t": "uintptr",
    "uint8_t": "uint8",
    "uint16_t": "uint16",
    "uint32_t": "uint32",
    "uint64_t": "uint64",
    "int32_t": "int32",
    "int64_t": "int64",
    "int": "int",
    "lore_node_id_t": "uint32",
    "void*": "unsafe.Pointer",
}


SEED_C_TYPES_THAT_NEED_CUSTOM_BUILDER = {
    "uint8_t": "bool",
    "lore_string_t": "string",
    "lore_string_array_t": "[]string",
    "lore_uint8_array_t": "[]bool",
    "lore_binary_t": "LoreBinary",
    "lore_bytes_t": "LoreBytes",
}


def detect_array_types(struct_dict, custom_types_set, use_ffi_map):
    """Find every `*_array_t` in `struct_dict` (excluding entries in
    `custom_types_set`) and return a list of dicts ready for Jinja loops.

    `use_ffi_map=True` mirrors `types.ji`'s element-resolution chain
    (cToGoFFITypeMap → cToGoTypeMap → pascal_case fallback). `False`
    mirrors `events.ji`'s chain, which historically skipped
    cToGoFFITypeMap. Preserving the asymmetry keeps generator output
    behaviorally identical.
    """
    arrays = []
    for c_type, fields in struct_dict.items():
        if not c_type.endswith("_array_t"):
            continue
        if c_type in custom_types_set:
            continue
        ptr_field = next((f for f in fields if f[1] == "ptr"), None)
        if ptr_field is None:
            continue
        ptr_type = ptr_field[0]
        if not ptr_type.endswith("*"):
            continue
        element_c_type = ptr_type.removesuffix("*")
        if use_ffi_map and element_c_type in SEED_C_TO_GO_FFI_TYPE_MAP:
            element_field_type = SEED_C_TO_GO_FFI_TYPE_MAP[element_c_type]
        elif element_c_type in SEED_C_TO_GO_TYPE_MAP:
            element_field_type = SEED_C_TO_GO_TYPE_MAP[element_c_type]
        elif element_c_type.endswith("_t"):
            element_field_type = util.pascal_case(element_c_type.removesuffix("_t"))
        else:
            element_field_type = element_c_type
        array_class = util.pascal_case(c_type.removesuffix("_t"))
        arrays.append(
            {
                "c_type": c_type,
                "array_class": array_class,
                "ffi_class": array_class + "FFI",
                "ptr_field_pascal": util.pascal_case(ptr_field[1]),
                "element_c_type": element_c_type,
                "element_field_type": element_field_type,
                "needs_element_builder": (
                    element_c_type in struct_dict
                    and element_c_type not in custom_types_set
                ),
            }
        )
    return arrays


def build_augmented(visitor):
    """Compose the dict of values to splat into Jinja globals.

    Templates access these by their bare name (e.g. `cToGoTypeMap`,
    `array_types`) — no `utils.` prefix.
    """
    custom_types_set = set(SEED_CUSTOM_TYPES)
    return {
        "ignored_functions": SEED_IGNORED_FUNCTIONS,
        "custom_types": SEED_CUSTOM_TYPES,
        "cToGoTypeMap": SEED_C_TO_GO_TYPE_MAP,
        "cToGoFFITypeMap": SEED_C_TO_GO_FFI_TYPE_MAP,
        "cTypesThatNeedCustomBuilder": SEED_C_TYPES_THAT_NEED_CUSTOM_BUILDER,
        "array_types": detect_array_types(
            visitor.types, custom_types_set, use_ffi_map=True
        ),
        "event_array_types": detect_array_types(
            visitor.events, custom_types_set, use_ffi_map=False
        ),
    }
