# CSVs CVR files (not all are vulnerable)

curl -C - -o v5-04-17-05_nm-cibola.csv https://ordros.com/cvr/New%20Mexico/Cibola/Cibola%202020%20CVR.csv
curl -C - -o v5-05-32-04_tn-williamson.csv https://ordros.com/cvr/Tennessee/Williamson/WilliamsonTNCVR.csv
curl -C - -o v5-05-12-01_ga-walker.csv https://ordros.com/cvr/2022Primaries/Georgia/CVR_Export_20220811140530_MAY2022.csv
curl -C - -o v5-05-52-06_ak.zip https://www.elections.alaska.gov/results/22SSPG/CVR_Export_20220908084311.zip
curl -C - -o v5-10-50-85_ca-placer.csv https://ordros.com/cvr/2022Primaries/California/CVR_Export_20220714101608_Placer.csv
curl -C - -o v5-11-03-01_co-pitkin.csv https://ordros.com/cvr/Colorado/Pitkin/cvr.csv
curl -C - -o v5-11-03-01_co-pitkin-2.csv https://ordros.com/cvr/Colorado/Pitkin/Pitkin%202020%20Redacted_Modified_CVR_Export_20201113121842.csv
curl -C - -o v5-12-05-39_nv-douglas.csv https://ordros.com/cvr/2022Primaries/Nevada/Nevada_Douglas_County_Primary_2022_CVR_Export_FOURTH_20220816161245.csv
curl -C - -o v5-15-16-01_nj-cumberland.csv https://ordros.com/cvr/2022Primaries/Cumberland/CVR_Export_20220303102155Cumberland.csv

# JSON-zip CVR Files (not all are vulnerable)

curl -C - -o v5-04-17-05_nm-otero.zip https://ordros.com/cvr/New%20Mexico/Otero/CVR_Export_20220120105928.zip
curl -C - -o v5-05-12-01_mi-wayne.zip https://ordros.com/cvr/Michigan/Wayne/CVR_Export_20220419121406.zip
curl -C - -o v5-05-32-04_az-maricopa.zip https://ordros.com/cvr/Arizona/Maricopa/CVR_Export_20210115215132.zip
curl -C - -o v5-05-40-02_ia-dickenson.zip https://ordros.com/cvr/Iowa/Dickenson/CVR_Export_20220824082504.zip
curl -C - -o v5-10-50-85_ca-san-fran.zip https://ordros.com/cvr/California/San%20Francisco/CVR_Export_20201201091840.zip
curl -C - -o v5-12-05-39_nv-douglas.zip https://ordros.com/cvr/2022Primaries/Nevada/Nevada_Douglas_County_Primary_2022_CVR_Export_FIRST_20220816154915.zip

# Ballot images

curl -C - -o ballot-images-1.7z https://www.zebraduck.org/election-files/georgia2021-01-05/ballot-images/chattahoochee-january-5.7z
curl -C - -o ballot-images-2.zip https://ordros.com/cvr/New%20Mexico/Chavez/Ballot%20Images.zip

# Check data downloaded is correct

sha256sum -c SHA256SUMS || exit 1

# Unpack ballot image data

if ! command -v 7z &> /dev/null
then
	@echo "********************* WARNING WARNING WARNING *********************"
	@echo "The 7z command is not available in the shell and the ballot-image"
	@echo "validation data can not be extracted."
	@echo ""
	@echo "You must decompress it manually to the correct location prior to"
	@echo "running validate.py."
	@echo "********************* WARNING WARNING WARNING *********************"
else
	7z x ballot-images-1.7z -o./ballot-images-1/
fi

unzip ballot-images-2.zip -d ballot-images-2
