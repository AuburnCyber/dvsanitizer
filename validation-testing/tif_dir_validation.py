import os
import hashlib
import re

def __load_directory(dir_path):
    hash_to_path_dict = {}
    for root, _, files in os.walk(dir_path):
        for a_file in files:
            file_path = os.path.join(root, a_file)
            with open(file_path, 'rb') as handle:
                file_hash = hashlib.sha256(handle.read()).hexdigest()

            # This is not a guarantee due to arbitrary directory structures but
            # is true for our testdata sources and makes validation
            # significantly simpler.
            assert file_hash not in hash_to_path_dict
            hash_to_path_dict[file_hash] = file_path

    return hash_to_path_dict

def compare_tif_dirs(dirty_base_dir, clean_base_dir, should_be_ignored_hashes):
    dirty_hash_to_path = __load_directory(dirty_base_dir)
    clean_hash_to_path = __load_directory(clean_base_dir)

    # CHECK: No files are unaccounted for
    num_clean = len(clean_hash_to_path) + len(should_be_ignored_hashes) - 1
    # files that were sanitized + files that should not be sanitized - skipped-during-sanitization.csv
    assert len(dirty_hash_to_path) == num_clean, '%d --- %d' % (len(dirty_hash_to_path), num_clean)

    skipped_filepath = None
    for clean_hash, clean_path in clean_hash_to_path.items():
        if clean_path.endswith('/skipped-during-sanitization.csv'):
            skipped_filepath = clean_path
            continue # Will handle afterwards

        # CHECK: No files added/changed
        assert clean_hash in dirty_hash_to_path

        dirty_path = dirty_hash_to_path[clean_hash]
        dirty_dir, dirty_filename = os.path.split(dirty_path)
        dirty_rel_dir = dirty_dir.removeprefix(dirty_base_dir)
        dirty_name, dirty_ext = os.path.splitext(dirty_filename)

        clean_dir, clean_filename = os.path.split(clean_path)
        clean_rel_dir = clean_dir.removeprefix(clean_base_dir)
        clean_name, clean_ext = os.path.splitext(clean_filename)

        # CHECK: Relative directory location is unchanged
        assert dirty_rel_dir == clean_rel_dir, dirty_rel_dir + ' --- ' + clean_rel_dir

        # CHECK: Extension is unchanged
        assert dirty_ext == clean_ext

        # CHECK: Non-record ID parts of filename are unchanged
        clean_name_split = clean_name.split('_')
        dirty_name_split = dirty_name.split('_')
        assert len(dirty_name_split) == 3 or len(dirty_name_split) == 4
        assert len(clean_name_split) == 3 or len(clean_name_split) == 4
        assert len(clean_name_split) == len(dirty_name_split)
        assert dirty_name_split[0] == clean_name_split[0], '%s --- %s' % (dirty_name_split[0], clean_name_split[0])
        assert dirty_name_split[1] == clean_name_split[1]
        if len(clean_name_split) == 4:
            assert dirty_name_split[3] == clean_name_split[3]

        clean_record_id = clean_name_split[2]
        dirty_record_id = dirty_name_split[2]

        # CHECK: Record ID was sanitized
        assert clean_record_id != dirty_record_id

        # CHECK: Record ID is expected format
        assert clean_record_id.startswith('0x')
        assert len(clean_record_id) == 18
        assert re.match('^0x[0-9a-f]{16}$', clean_record_id)

        del dirty_hash_to_path[clean_hash]

    # CHECK: A skipped-file was found
    assert skipped_filepath is not None

    # CHECK: All of the should-be-skipped files were skipped
    for should_skip_hash in should_be_ignored_hashes:
        assert should_skip_hash in dirty_hash_to_path
        del dirty_hash_to_path[should_skip_hash]

    # CHECK: There are no unexpected files left-over
    assert len(dirty_hash_to_path) == 0
