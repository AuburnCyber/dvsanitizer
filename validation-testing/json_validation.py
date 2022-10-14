import json
import os
import re
import zipfile

def __compare_inner_cvr_files(dirty_zip_handle, clean_zip_handle, filename):
    with dirty_zip_handle.open(filename) as handle:
        dirty_obj = json.loads(handle.read())
    with clean_zip_handle.open(filename) as handle:
        clean_obj = json.loads(handle.read())

    # CHECK: Root-level object has only expected pairs
    assert 'Version' in clean_obj
    assert 'Sessions' in clean_obj
    assert 'ElectionId' in clean_obj
    assert len(clean_obj.keys()) == 3

    # CHECK: SW version and Election Id is unchanged
    assert clean_obj['Version'] == dirty_obj['Version']
    assert clean_obj['ElectionId'] == dirty_obj['ElectionId']

    # CHECK: The location for session objects exists
    assert 'Sessions' in clean_obj

    clean_record_ids = set()
    clean_reduced_sessions = []
    for clean_session in clean_obj['Sessions']:
        # CHECK: Record ID exists and is the expected format
        assert 'RecordId' in clean_session
        assert type(clean_session['RecordId'] is str)
        assert re.match('^0x[0-9a-f]{16}$', clean_session['RecordId'])
        assert len(clean_session['RecordId']) == 18


        # CHECK: ImageMask exists and was updated along with record ID
        record_id_str = clean_session['RecordId']
        assert 'ImageMask' in clean_session
        _, name_glob = os.path.split(clean_session['ImageMask'])
        split_name_glob = name_glob.split('_')
        assert split_name_glob[2].startswith(record_id_str) # Length 3 and 4 w/o hassle from star

        # CHECK: All record IDs are unique
        assert clean_session['RecordId'] not in clean_record_ids
        clean_record_ids.add(clean_session['RecordId'])

        # Though each session is complex with many nested objects, cheaply
        # compare using the 'rest-of-line' approach where everything except the
        # record ID and image mask should be the same. Passing 'sort_keys=True'
        # ensures deterministic in marshalling to a string.
        del clean_session['RecordId']
        del clean_session['ImageMask']
        clean_as_str = json.dumps(clean_session, sort_keys=True)
        clean_reduced_sessions.append(clean_as_str)

    # Just walk, checks are afterwards for simplicity
    dirty_record_ids = [] # *not* guaranteed to be unique
    dirty_reduced_sessions = []
    for dirty_session in dirty_obj['Sessions']:
        dirty_record_ids.append('0x%06x' % dirty_session['RecordId']) # Make easy to look for dupes
        del dirty_session['RecordId']
        del dirty_session['ImageMask']
        dirty_as_str = json.dumps(dirty_session, sort_keys=True)
        dirty_reduced_sessions.append(dirty_as_str)

    # CHECK: Same number of sessions in both
    assert len(dirty_record_ids) == len(clean_record_ids)
    assert len(dirty_reduced_sessions) == len(clean_reduced_sessions)
    assert len(clean_record_ids) == len(clean_reduced_sessions)
    num_sessions = len(clean_record_ids)

    # CHECK: All record IDs were sanitized
    assert len(clean_record_ids.intersection(set(dirty_record_ids))) == 0

    # CHECK: Nothing else changed unexpectedly in the sessions
    sorted_dirty = sorted(dirty_reduced_sessions)
    sorted_clean = sorted(clean_reduced_sessions)
    for i in range(num_sessions):
        assert sorted_dirty[i] == sorted_clean[i]

    return clean_record_ids

def compare_json_zip_files(dirty_path, clean_path):
    dirty_zip_handle = zipfile.ZipFile(dirty_path, mode='r')
    dirty_filenames = dirty_zip_handle.namelist()
    clean_zip_handle = zipfile.ZipFile(clean_path, mode='r')
    clean_filenames = clean_zip_handle.namelist()

    # CHECK: No files were added or lost during sanitizing
    assert sorted(dirty_filenames) == sorted(clean_filenames)
    to_be_processed = clean_filenames

    all_clean_record_ids = set()
    for filename in to_be_processed:
        print('Comparing internal zip-file: ' + filename)
        if not filename.startswith('CvrExport') or not filename.endswith('.json'):
            # CHECK: Non-CVR files are unchanged
            with dirty_zip_handle.open(filename) as handle:
                dirty_data = handle.read()
            with clean_zip_handle.open(filename) as handle:
                clean_data = handle.read()
            assert dirty_data == clean_data

            continue

        new_record_ids = __compare_inner_cvr_files(dirty_zip_handle, clean_zip_handle, filename)

        # CHECK: There are no cross-file record ID collisions
        assert len(all_clean_record_ids.intersection(new_record_ids)) == 0
        all_clean_record_ids = all_clean_record_ids.union(new_record_ids)

    dirty_zip_handle.close()
    clean_zip_handle.close()
