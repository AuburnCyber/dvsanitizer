import re

def __split_csv_line(complete_line, is_clean):
    split_line = complete_line.split(',')

    record_id = split_line[3]
    del split_line[3] # remove from rest-of-line b/c is intentionally not the same

    cvr_id = split_line[0]
    del split_line[0] # remove from rest-of-line b/c intentionally re-ordering makes not consistent

    rest_of_line = ','.join(split_line) # data line without record or CVR ID

    # IDs can be encoded in multiple ways depending on many factors.
    if is_clean:
        assert record_id != '='

        if record_id.startswith('"'):
            assert record_id[0] == '"'
            assert record_id[-1] == '"'
            record_id = record_id[1:-1]

            is_oddball_encoding = True

        else:
            assert record_id[0] != '='
            assert record_id[0] != '"'
            assert record_id[-1] != '"'

            is_oddball_encoding = False
    else:
        # Dirty record IDs and CVR IDs can be as a raw number (',123,') or as
        # in the oddball quoting format (',="123",').
        if record_id.startswith('="'):
            assert record_id[0] == '='
            assert record_id[1] == '"'
            assert record_id[-1] == '"'
            record_id = record_id[2:-1]

            is_oddball_encoding = True
        else:
            is_oddball_encoding = False

    if is_clean:
        # CHECK: Record ID has the expected format
        assert re.match('^0x[0-9a-f]{16}$', record_id), record_id
        assert len(record_id) == 18

        # CHECK: CVR ID is deleted
        if is_oddball_encoding:
            assert cvr_id[0] == '=', cvr_id
            assert cvr_id[1] == '"', cvr_id
            assert cvr_id[-1] == '"'
            assert len(cvr_id) == 3
        else:
            assert len(cvr_id) == 0, cvr_id
    else:
        assert record_id.isnumeric()


    return record_id, rest_of_line, is_oddball_encoding

def __load_clean_csv_file(clean_path):
    with open(clean_path, 'r') as clean_handle:
        clean_data = clean_handle.read()
    all_clean_lines = clean_data.splitlines()
    clean_header_lines = all_clean_lines[:4]
    clean_data_lines = all_clean_lines[4:]

    # Load data to check against dirty-data and validate that all clean CVR IDs are unique.
    lookup_dict = {} # rest-of-line -> count
    clean_record_ids = []
    oddball_encoded = None
    for line in clean_data_lines:
        record_id, clean_rol, is_oddball = __split_csv_line(line, True)
        if oddball_encoded is None: # First line parsed
            oddball_encoded = is_oddball
        else:
            # CHECK: all lines are either oddball or not
            assert oddball_encoded == is_oddball

        # CHECK: Record ID is unique in this file.
        assert record_id not in clean_record_ids
        clean_record_ids.append(record_id)

        # Add rest-of-line to the tracking dict OR update the exisiting counter
        if clean_rol not in lookup_dict:
            lookup_dict[clean_rol] = 0
        lookup_dict[clean_rol] += 1

    return clean_header_lines, clean_record_ids, lookup_dict, oddball_encoded

def __load_dirty_csv_file(dirty_path):
    with open(dirty_path, 'r') as dirty_handle:
        dirty_data = dirty_handle.read()
    all_dirty_lines = dirty_data.splitlines()
    dirty_header_lines = all_dirty_lines[:4]
    dirty_data_lines = all_dirty_lines[4:]

    # Load data but don't need to check anything about this data.
    lookup_dict = {} # rest-of-line -> count
    dirty_record_ids = []
    oddball_encoded = None
    for line in dirty_data_lines:
        record_id, dirty_rol, is_oddball = __split_csv_line(line, False)
        if oddball_encoded is None: # First line parsed
            oddball_encoded = is_oddball
        else:
            assert oddball_encoded == is_oddball

        dirty_record_ids.append(record_id)

        if dirty_rol not in lookup_dict:
            lookup_dict[dirty_rol] = 0
        lookup_dict[dirty_rol] += 1

    return dirty_header_lines, dirty_record_ids, lookup_dict, oddball_encoded

def compare_csv_files(dirty_path, clean_path):
    clean_header_lines, clean_record_ids, clean_rol_dict, clean_is_oddball = __load_clean_csv_file(clean_path)
    dirty_header_lines, dirty_record_ids, dirty_rol_dict, dirty_is_oddball = __load_dirty_csv_file(dirty_path)

    # CHECK: Both files are encoded in the same way
    assert clean_is_oddball == dirty_is_oddball

    # CHECK: The header lines are the unchanged
    assert len(dirty_header_lines) == 4
    assert len(clean_header_lines) == 4
    assert dirty_header_lines[0] == clean_header_lines[0]
    assert dirty_header_lines[1] == clean_header_lines[1]
    assert dirty_header_lines[2] == clean_header_lines[2]
    assert dirty_header_lines[3] == clean_header_lines[3]

    # CHECK: Clean record IDs are unique
    assert len(set(clean_record_ids)) == len(clean_record_ids)

    # CHECK: There are no unchanged record IDs (impossible due to '0x' prefix)
    assert len(set(dirty_record_ids).intersection(set(clean_record_ids))) == 0

    # CHECK: No lines were lost
    num_clean_data_lines = 0
    for line_count in clean_rol_dict.values():
        num_clean_data_lines += line_count
    num_dirty_data_lines = 0
    for line_count in dirty_rol_dict.values():
        num_dirty_data_lines += line_count
    assert num_clean_data_lines == len(clean_record_ids)
    assert num_clean_data_lines == num_dirty_data_lines

    # CHECK: ROL uniqueness is unchanged
    assert len(clean_rol_dict) == len(dirty_rol_dict)
    assert sorted(clean_rol_dict.keys()) == sorted(dirty_rol_dict.keys())

    # CHECK: Unique ROL counts are unchanged
    for clean_rol, clean_count in clean_rol_dict.items():
        assert clean_rol in dirty_rol_dict
        assert clean_count == dirty_rol_dict[clean_rol]
