import argparse
import hashlib
import os
import shutil
import subprocess
import sys

from data_for_validation import get_tif_dir_ignored, get_expected_sanitized_hash
from csv_validation import compare_csv_files
from json_validation import compare_json_zip_files
from tif_dir_validation import compare_tif_dirs

APP_PATH = '../dvsanitizer'
INSECURE_SEED_FOR_TESTING = 'aaaaaaaaaaaaaaaa'
LARGE_ZIP_THRESHOLD = 200000000 # in bytes (currently 200MB)
TESTDATA_FILES_TO_IGNORE = set([
    'download.sh',
    'SHA256SUMS',
    'ballot-images-1.7z',
    'ballot-images-2.zip',
    '.gitignore',
    ])

arg_parser = argparse.ArgumentParser()
arg_parser.add_argument('--include-large-zips',
                        default=False,
                        dest='include_large_zips',
                        action='store_true',
                        help='Validate large zip-files (>200MB) instead of refusing',
                        )
arg_parser.add_argument('--skip-large-zips',
                        default=False,
                        dest='skip_large_zips',
                        action='store_true',
                        help='Skip large zip-files (>200MB) instead of refusing',
                        )
args = arg_parser.parse_args()
if args.include_large_zips and args.skip_large_zips:
    sys.exit('Impossible to allow and skip large JSON-Zips')


def __build_binary():
    cmd = ['go', 'build']
    print('CMD: ' + ' '.join(cmd))
    subprocess.run(cmd, check=True, cwd='../')

def __run_csv_file(test_file):
    if os.path.isfile('./output/' + test_file):
        os.remove('output/'+test_file)
        print('Removed existing output file')

    cmd = [
            APP_PATH,
            '--sanitize-csv',
            '--seed', INSECURE_SEED_FOR_TESTING,
            '--output-dir', 'output/',
            '--input', os.path.join('testdata', test_file),
            ]
    print('CMD: ' + ' '.join(cmd))
    subprocess.run(cmd, check=True)

    compare_csv_files(
            os.path.join(os.getcwd(), 'testdata', test_file),
            os.path.join(os.getcwd(), 'output', test_file),
            )
    print('CSV validation passed')

def __run_json_zip_file(test_file):
    if os.path.isfile('./output/' + test_file):
        os.remove('output/'+test_file)
        print('Removed existing output file')


    cmd = [
            APP_PATH,
            '--sanitize-json-zip',
            '--seed', INSECURE_SEED_FOR_TESTING,
            '--output-dir', 'output/',
            '--input', os.path.join('testdata', test_file),
            ]
    print('CMD: ' + ' '.join(cmd))
    subprocess.run(cmd, check=True)

    compare_json_zip_files(
            os.path.join(os.getcwd(), 'testdata', test_file),
            os.path.join(os.getcwd(), 'output', test_file),
            )
    print('JSON-Zip validation passed')

def __run_tif_sha_dir(test_dir):
    if os.path.isdir('./output/' + test_dir):
        shutil.rmtree('output/'+test_dir)
        print('Removed existing output directory')

    os.mkdir('output/'+test_dir)
    print('Created new output directory')

    cmd = [
            APP_PATH,
            '--sanitize-tif-dir',
            '--seed', INSECURE_SEED_FOR_TESTING,
            '--output-dir', os.path.join('output', test_dir),
            '--input', os.path.join('testdata', test_dir),
            ]
    print('CMD: ' + ' '.join(cmd))
    subprocess.run(cmd, check=True)

    compare_tif_dirs(
            os.path.join(os.getcwd(), 'testdata', test_dir),
            os.path.join(os.getcwd(), 'output', test_dir),
            get_tif_dir_ignored(test_dir),
            )
    print('TIF-dir validation passed')

def main():
    __build_binary()

    validated_csvs = []
    validated_json_zips = []
    validated_tif_dirs = []
    for test_file in os.listdir('testdata'):
        if test_file in TESTDATA_FILES_TO_IGNORE:
            # Silently ignore b/c expected to be there and not meant to be sanitized.
            continue

        should_check_determinism = False # Only for CSV and ZIP files
        if test_file.endswith('.csv'):
            print('Validating CSV: ' + test_file)
            __run_csv_file(test_file)
            should_check_determinism = True
            validated_csvs.append(test_file)

        elif test_file.endswith('.zip'):
            if os.path.getsize(os.path.join('testdata', test_file)) > LARGE_ZIP_THRESHOLD:
                if args.skip_large_zips:
                    print('Skipping large JSON-Zip: ' + test_file)
                    continue
                if not args.include_large_zips:
                    print('-- FAIL --' * 10)
                    print('-' * 80)
                    print('Large JSON-Zip encountered: ' + test_file)
                    print('Cowardly refusing to validate without --include-large-zips.')
                    print('This could take many hours to validated depending on hardware.')
                    print('-' * 80)
                    print('-- FAIL --' * 10)
                    sys.exit(1)

            print('Validating JSON-Zip: ' + test_file)
            __run_json_zip_file(test_file)
            should_check_determinism = True
            validated_json_zips.append(test_file)

        elif test_file in ['ballot-images-1', 'ballot-images-2']:
            print('Validating tif-dir:' + test_file)
            __run_tif_sha_dir(test_file)
            validated_tif_dirs.append(test_file)

        else:
            print('IGNORE UNKNOWN FILE: ' + test_file)
            continue

        # CHECK: Determinism
        if should_check_determinism:
            # Avoid hard-to-debug error b/c changes expected hashes
            assert INSECURE_SEED_FOR_TESTING == 'aaaaaaaaaaaaaaaa'

            want_hash = get_expected_sanitized_hash(test_file)
            with open('output/'+test_file, 'rb') as handle:
                found_hash = hashlib.sha256(handle.read()).hexdigest()
            assert found_hash == want_hash
            print('Validated determinism:' + test_file)
        else:
            print('Ignoring determinism:' + test_file)

        print('Finished validating: ' + test_file)

    print('FINISHED SUCCESSFULLY')
    print('CSVs validated: ' + str(validated_csvs))
    print('JSON-Zips validated: ' + str(validated_json_zips))
    print('TIF-dirs validated: ' + str(validated_tif_dirs))

main()
