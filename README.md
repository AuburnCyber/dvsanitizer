# DVSanitizer

DVSanitizer is an automated tool to assist election officials in reprocessing Dominion
cast-vote records (CVRs) and ballot image files to prevent exploitation of the
[DVSorder privacy flaw](https://DVSorder.org), which affects data from ICP and ICE
tabulators. The tool will no longer be needed after Dominion effectively patches
the the flaw.

Sanitizing ballot-level data _cannot_ affect official election results, because results are generated directly from the election management system (EMS), not from the ballot-level data released to the public. However, **as with any third-party software, jurisdictions should not run this sanitization tool on their EMS computers**. Doing so may violate certification conditions applicable to the EMS. Instead, we recommend copying vulnerable CVRs or ballot images to an external system and running the tool there. The tool is open-source software, and we encourage anyone interested to view the code and test its behavior.

This software was created by researchers at Auburn University and the University of Michigan.
We'd be happy to provide election officials any assistance we can in using this tool or sanitizing data.
Please reach out at [team@DVSorder.org](mailto:team@DVSorder.org).

## How It Works

DVSanitizer can process CVRs in .csv or .zip format and folders of ballots images in .tif format.                                                           

It works by replacing all instances of ballot record IDs with AES-encrypted versions of those IDs,
using a user-supplied or randomly generated secret key.
Once this key is destroyed or forgotten, it is infeasible to unshuffle the ballots.

Encrypting the record IDs has two advantages over simply deleting them:

1. By using the same key to sanitize CVRs and ballot images, officials can ensure that record IDs in both are encrypted the same way. This preserves the ability to cross-reference the CVR entries and image files after sanitization.

2. If CVRs or ballot images are published incrementally during the counting process, officials can use the same key to sanitize all data releases, thus ensuring that the sanitized record IDs remain consistent. This preserves the ability to identify which ballots have been added or changed across the different releases.

## Basic Examples

#### To sanitize a single CSV-format CVR file

```
dvsanitizer --gen-seed --sanitize-csv --input dirty-data/CVR_Export_1234.csv --output-dir clean-data/
```

Data will be read from `dirty-data/CVR_Export_1234.csv` and the sanitized version
will be written to `clean-data/CVR_Export_1234.csv`.

**Note that the input data files are not modified.** The data is read, processed, and then
written to new files in the output directory, leaving the input files unchanged.


#### Sanitize a single JSON-format CVR .zip file:

```
dvsanitizer --gen-seed --sanitize-json-zip --input dirty-data/CVR_Export_1234.zip --output-dir clean-data/
```

#### To sanitize a folder (or folder hierarchy) of .tif ballot images:

```
dvsanitizer --gen-seed --sanitize-tif-dir --input dirty-data/ballot-images/ --output-dir clean-data/ballot-images/
```

## Advanced Example

If you plan to publish multiple kinds of ballot-level data from an election,
or multiple versions of the same kind of data from an election,
it is helpful to sanitize them all using the same seed.
This preserves the ability to cross-reference
the same ballot between sanitized datasets.

To do so, first run the tool with the `--gen-seed` flag:

```
dvsanitizer --gen-seed --sanitize-csv --input dirty-data/CVR_Export_1234.csv --output-dir clean-data/
```

When the process complete, the tool will dislay the randomly generated seed it used (replaced with `XXX`s below):

```
********************************************************************************
********************************************************************************
Your auto-generated seed is: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
********************************************************************************
********************************************************************************
```

You can provide this seed when sanitizing subsequent files using the `--seed` flag.

For example, to sanitize ballot images so that the use IDs that are consistent with the CVR file,
you could run this (replacing the `XXX`s with the seed value printed by the first command):

```
dvsanitizer --seed=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --sanitize-tif-dir --input dirty-data/ballot-images/ --output-dir clean-data/ballot-images/
```


## Usage Details

```
dvsanitizer
  [--gen-seed OR --seed CHOSEN-SEED]
  [--sanitize-csv OR --sanitize-json-zip OR --sanitize-tif-dir]
  --input [FILENAME or TIF-DIRECTORY]
  --output-dir DIRECTORY
```

In order, these options are:

* Seed Instruction (`--gen-seed` or `--seed`)
    * If `--gen-seed` is passed, a cryptographically secure value is generated,
      used for this instance, and then printed to the screen once sanitization
      is complete
    * If `--seed` is passed, the given value is used to encrypt the IDs
    * **Using the same seed on the same or overlapping sets data ensures sanitized record IDs are matched across files/formats**
* Action (`--sanitize-csv`, `--sanitize-json-zip`, or `--sanitize-tif-dir`)
    * Tells DVSanitizer the format of the data to be sanitized
      to correctly sanitize that data format
    * A couple of important things to note about using `--sanitize-tif-dir`:
        * Will recursively look for `.tif` and `.sha` files under that folder
          and sanitize any it finds
        * The directory structure will be re-created in the output location to
          the extent is required to write the sanitized files
        * A file named `skipped-during-sanitization.csv` will be added and will
          contain a list of files seen but not copied to the output location 
          as well as a reason why each was skipped
* Input Location (`--input`)
    * Tells DVSanitizer which file or folder you want to sanitize
    * Must be a folder if the action is `--sanitize-tif-dir` or a file for the other actions
* Output Directory (`--output-dir`)
    * Tells DVSanitizer where to write the output
    * DVSanitizer will pick the output filename(s) within this folder based on the input data
        * The filenames for CVR files remain the same after
          sanitization
        * The filenames for ballot images change due to the
          inclusion of the record ID in the filename, but the directory
          structure is unchanged
    * The command will fail if an output file exists. The program **will not** overwrite 
      existing files
    * If the action is `--sanitize-tif-dir`, the output directory **must** be
      empty.

## Data Changes During Sanitization

DVSanitizer attempts to avoid altering more data than necessary.
Below is a list of all intentional data changes to the data during sanitization
as well as some notable data aspects that *are not* changed.

#### Complete Listing of Modified Data

| Format | Field | Rationale |
| :----- | :---- | :-------------- |
| CSV | `RecordId` | Replaced with sanitized record ID to remove ability to directly deanonymize ballots |
| CSV | *row order* | Re-ordered to avoid side-channel from unsanitized record ID ordering |
| CSV | `CvrNumber` | Removed to avoid side-channel from unsanitized record ID ordering |
| JSON | `Sessions[].RecordId` | Replaced with sanitized record ID to remove ability to directly deanonymize ballots |
| JSON | `Sessions[].ImageMask` | Replaced portion of filename with sanitized record ID to remove ability to directly deanonymize ballots |
| JSON | `Sessions[]` | Re-ordered to avoid side-channel from unsanitized record ID ordering |
| TIF/SHA | *filename* | Replace portion of filename with sanitized record ID to remove ability to directly deanonymize ballots |

#### Notable Unmodified Data

* CVR elements representing voter's selections in either the CSV or JSON format
* Any data within the `.tif` or `.sha` file (only filename modified)
* `TabulatorId`/`TabulatorNumber` or `BatchId` in CSV or JSON format
* Any description of the precinct, ballot type, or voting method

## Legal Disclaimers

ALL CONTENT, SERVICES, PRODUCTS, AND SOFTWARE - SUCH AS THE DVSanitizer TOOL - PROVIDED IN THIS REPOSITORY OR BY DVSorder.org ARE PROVIDED 'AS IS' WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED. WE DISCLAIM ALL WARRANTIES, EXPRESS OR IMPLIED, INCLUDING, WITHOUT LIMITATION, THOSE OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NONINFRINGEMENT. YOU ARE SOLELY RESPONSIBLE FOR THE APPROPRIATENESS AND CONSEQUENCES OF YOUR USE OF THE DVSanitizer TOOL, AS WELL AS ANY OTHER CONTENT, SERVICES, PRODUCTS, AND SOFTWARE IN THIS REPOSITORY OR BY DVSorder.org.  WE DO NOT WARRANT THAT THAT THE CONTENT, SERVICES, PRODUCTS, AND SOFTWARE PROVIDED IN THIS REPOSITORY OR BY DVSorder.org MEET YOUR REQUIREMENTS.  BY USING ANY CONTENT, SERVICES, PRODUCTS, OR SOFTWARE IN THIS REPOSITORY OR BY DVSorder.org YOU UNDERSTAND AND AGREE THAT WE SHALL NOT BE LIABLE FOR ANY DIRECT, INDIRECT, SPECIAL, CONSEQUENTIAL, INCIDENTAL, OR PUNITIVE DAMAGES.  YOU USE ANY CONTENT, SERVICES, PRODUCTS OR SOFTWARE ON THIS SITE AT YOUR OWN DISCRETION. As with any third-party software, **jurisdictions should not run DVSanitizer on their EMS computers**.
