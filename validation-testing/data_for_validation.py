# These are easy to generate for the non-tif and -non-sha files when adding new
# validation datasets the below command but have to add individual hashes for
# things like 'NotCast'.
#
#       find <<<dir>>> -type f -print0 | xargs -0 sha256sum | grep -v '\.tif$' | grep -v '\.sha$' | sed 's/  /", # /' | sed 's/^/"/'
def get_tif_dir_ignored(dir_name):
    if dir_name == 'ballot-images-1':
        return [
                "f53ff9c17df8250ee4d2f040c23f648048cd8ce2a817df392e7937493d1beec4", # ballot-images-1/January 2020 Runoff/Tabulator00060/Results/1_20_60_1_DETAIL.DVD
                "3f594b22956967dd79cd841f62160d06a916c51a27d5b833fb66935a64b31d99", # ballot-images-1/January 2020 Runoff/Tabulator00060/Results/1_20_60_1_TOTALS.DVD
                "d2f2288a14d43b2a2dd5fd619cd219116740683e953d276a65a94ac2ae5272ec", # ballot-images-1/January 2020 Runoff/Tabulator00060/Results/1_20_60_1_RAW.DVD
                "8552614ba0db36d4d3c7c6c3f420805bdd2fbe8630020a2059cdf536b8759281", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_9_RAW.DVD
                "1b60b95f0a5884b0c74365df384823ed332208d4230c167c1a551319d0791b86", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_8_TOTALS.DVD
                "7682b72ca0b145dcde170420d079e0ca5aef1921ee1f666c1ce8aefdc7bbad48", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_12_RAW.DVD
                "1c3b8f697431274628be9676a1e3e8d37f680b002b36a9c4c74b4fb7d4c930c2", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_5_TOTALS.DVD
                "aac50a2eee31c2335ac7c1755440e95f3fb3efc9067b476a0af741fc884cfffa", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_11_TOTALS.DVD
                "05b9e58a29b8c6d64a23b2ec8c0b0d76bbd2e13be9fe191d1a32ecb1788b7fc2", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_3_DETAIL.DVD
                "81c7c619019689aeac7b9ad518ad198c3ca418342162528da0f4d2397e76c791", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_17_DETAIL.DVD
                "f1e0cc44005d033d634dbf6785e756ea333d81020e783ef44cc9c43e6ac4f426", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_14_TOTALS.DVD
                "41aa92e1faf4cd0f0edbc4addd370d058cbb6a099b2303c6bb79abd21daed99c", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_1_RAW.DVD
                "38c0df190e094e5db9940ed1b074a78b963d8b6837286ee3f342678fd961c886", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_12_DETAIL.DVD
                "d9db5c6e54bf63ea69b9b5cbe72c2af30ac60423de323a068028e31f4ab4eb88", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_6_DETAIL.DVD
                "0bd24df205a61c891c6a73243ac01e513fc918e9772869eb0645f92723bbdcab", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_8_RAW.DVD
                "79e6cef99ad56f3a822f54521980e722813540d0a03d14fc25b473a5d32d3f7e", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_13_RAW.DVD
                "e33a0c8d5b492743fca4038272c1fde3ffd14692b75633d04de9dcb6f0a04621", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_3_RAW.DVD
                "18ea96354abaeb34bd6e7f78611db33375e243a2c48c28a7d9b11fe2258a855b", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_11_RAW.DVD
                "512bfd7352092883b403275d812d8a0b53e9d7344308c82291dc5cd5db9c40aa", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_1_DETAIL.DVD
                "2a402e16249ff2a5c37657e7f64a4ed260538b90ee9764b9eb3349fdd85d633a", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_15_DETAIL.DVD
                "2fc5f15484fa2b0334855d011b44ac2c8650e25fe3ddbb19d9460d727859689e", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_7_TOTALS.DVD
                "1f5c5c58f99f91ae5ec6248fa60b8a1ce2aab5025b90f095d299a8e6d5381904", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_13_TOTALS.DVD
                "a9291a6a9f5f4f17fb8ee92712690eab79738860fd5aefe7ddff32e15f084934", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_10_DETAIL.DVD
                "6399a6155044ae5473752b2f2abb9106a95147e1ded8c400183e2af56b802d1d", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_10_RAW.DVD
                "cdd86c18ff5a63c7ae49d1b8fd742e938991c2c1674ebf1f5cbb49241df0d713", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_4_DETAIL.DVD
                "072c54daa84e50d147ddcff04520e125852d082a93ab8040e3b0cee25bd621db", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_16_TOTALS.DVD
                "2df2580f8eaaf732e2072d92db6fb1dde11418935b67aa6f209fd33a3e81185c", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_2_TOTALS.DVD
                "852e32a65cc98161f792ec45686733d30f4e80360e657ae647899061c1a22a6e", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_9_DETAIL.DVD
                "a92a60ac76282e9a2d3720074527a8c754d8d15d2b5e8f2956de7b1f1411f0c8", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_2_RAW.DVD
                "2dd9c72ef56b52b1b79cbaed189e415eb320b616b00eaf2de29fd48b34c0e83f", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_1_TOTALS.DVD
                "87d444fdab1b15d4c646303f531de8f009400f4e8242bb171b0dbaa0cc07fa04", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_15_TOTALS.DVD
                "c1efc00476f1897b951759b7e14993d63befad4b7136c16e20523683573c92e9", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_7_RAW.DVD
                "6eda7e8631554735d2c1d6eb6d7e34c93c383d819cb6cc0d660b0e50f9ffeede", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_7_DETAIL.DVD
                "168fa8f9d2f4c1da86e0fcf4ae9cb74a80995439e451e39560d4b8ec276c3d26", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_13_DETAIL.DVD
                "40c64ad1fe21c007ea26740c204b31192bac21289e818c2aa628430337d3d9fd", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_15_RAW.DVD
                "267af8ee4a06c8dcd0739b95c424dc8c89fcb7cb067fa04250f6c34e33d9b5bf", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_9_TOTALS.DVD
                "02e33de6541a6f8cc9ff46f13ee696f47c90e7bc870c7c74cb34ee2c0114a733", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_14_RAW.DVD
                "8a530695c74d6fb3e4421ff63ec30a87c4acf8c8fec9745e4d9f57f5827e684b", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_10_TOTALS.DVD
                "7e79395fde05e3f5c83455fced89ba006e7d638647241643e06d6de24f9ffcc5", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_4_TOTALS.DVD
                "0e6b025d6655a20044a4d9d733d31b3a779caf4dbfad08ed056a0e7a7dba4ed2", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_6_RAW.DVD
                "869c59d8332f04d8e3a908a79ac3c3666051439af21d53f0781732845cbcf282", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_16_DETAIL.DVD
                "38cb67e682e19c74dbfc671d51cb52b3a4aa9137f287459c32256bfb8cda0edc", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_2_DETAIL.DVD
                "9465aafd4001d4ef3c67a1f7446ccf83a338969f67b46a5ca9cc5cb2b17b1c43", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_5_DETAIL.DVD
                "83ee4cd5811a1d7dd6a22c37a603565545a303952185821cfa86ca6a6ce21ed6", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_11_DETAIL.DVD
                "a21faa774bba5b7aa292152b9f3e843476222e65dba7b3b0761b4da4aa19b1ff", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_3_TOTALS.DVD
                "06757844cf72201214162702345790db40f39ea787ff8e229a7cbf118a617122", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_16_RAW.DVD
                "98e0cb7aaf0219f2d6a080c27b4bbcb316422caf7db4ad967d740967d9c458fd", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_17_TOTALS.DVD
                "bf82db5c7635de292d8efd411d1a0efdcd17eb6ae09d61d9d4a8c0a5520d0a08", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_8_DETAIL.DVD
                "e9bde8438b06e2d3e8c346fc4fb29e4474b5ac43a8e30a02f9ad3b7f15234735", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_4_RAW.DVD
                "632f6b7bc10cd94ac5700cd20d7647e7c87684a583418224f5777ac0445170eb", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_5_RAW.DVD
                "ffd0569a7c5e44c71a27b05c66d035da169cdab52fbcf14749be44a4adcd4122", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_14_DETAIL.DVD
                "14875c0115c70aa40dacd8bd7b6b3619d8afa77e2b60f3eca5dca7a12b87544f", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_12_TOTALS.DVD
                "fdacbed9c81d1ae71811891db22986adef5a878e1137fb203e845520785c7f6c", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_6_TOTALS.DVD
                "5495efe361332a1785ecf4b5ef12fde636bd915b41610fd3c96ce24d36029ac2", # ballot-images-1/January 2020 Runoff/Tabulator00050/Results/1_20_50_17_RAW.DVD
                "3c299735cf05def7c46e4d417fa6dbf8272baf45d3b9f148e3af761746c200d5", # ballot-images-1/Additional Documents/scan0002.pdf
                "b8d54ba1da1ca7d153c96d55c46da3ec1990747af4784e45c93ad36b3da7b297", # ballot-images-1/Additional Documents/scan0003.pdf
                "d2b24e5d53048042955074a2c23cca0e518b569cb6352e59f870dc920b52d489", # ballot-images-1/Additional Documents/scan0004.pdf
                "d5213dd6d39a9405688faeeaa5fe6d8f9c0327ec5e371c480ca49166551af60d", # ballot-images-1/Additional Documents/scan0005.pdf
                ]
    if dir_name == 'ballot-images-2':
        return [
                "96afa5c97e2015e5ec8512d6c7ec3734f0e6e9a15fcc7b090fb1be19d4498a01", # ballot-images-2/__2020 GENERAL - BALLOT IMAGES AND HAND TALLIES/Tabulator00983/Batch000/Images_Error.txt
                "aaf11b180613baa064a4f38c0de2ea1f701401c1f96cab6b382780735f985947", # ballot-images-2/__2020 GENERAL - BALLOT IMAGES AND HAND TALLIES/Tabulator00988/Batch000/Images_Error.txt
                "67b3e58f892f5c8375dc2df1b6975e938f4df118adefe86b9e08a4a990298667", # ballot-images-2/___All_Errors.txt
                "7c03ca6843b9e78875093cc5c336a31a1e7acfe23951678d3e78d76680f2428f", # ballot-images-2/2020 GENERAL - BALLOT IMAGES AND HAND TALLIES/Absentee Hand Tallies - Pct 73 - 106.pdf
                ]


    assert False, dir_name

def get_expected_sanitized_hash(filename):
    if filename == 'v5-11-03-01_co-pitkin.csv':
        return 'c671ba49248c5920cee641fc2eb8c941ec4269f0734feb70a4c250d75bbf0b8a'
    if filename == 'v5-05-40-02_ia-dickenson.zip':
        return '90d36f637d162d11243628d8e2ec87ea6603bc8c94e97c4f11c403fe713c5f71'
    if filename == 'v5-10-50-85_ca-san-fran.zip':
        return '3714245190f93d386c9188b2387af38f829a166196c2f09304c21ffb1772aa27'
    if filename == 'v5-04-17-05_nm-otero.zip':
        return '7f48c91bd5e260f9ad27eb157560bcf2f479fd6a8088335339a172ee0bdb083b'
    if filename == 'v5-12-05-39_nv-douglas.csv':
        return 'b017401541e4742bc2d82906bf7c7532cde74947d4206d3028d22ab8d01c70ac'
    if filename == 'v5-10-50-85_ca-placer.csv':
        return '117176ec6786a67aa570997d82d2c6e8f822ae426de6a4b4e4c6def511879b29'
    if filename == 'v5-12-05-39_nv-douglas.zip':
        return '8976759416adf277b1863fa48b436d1f73a6c365347d4177b0269d90d922aad8'
    if filename == 'v5-11-03-01_co-pitkin-2.csv':
        return 'c671ba49248c5920cee641fc2eb8c941ec4269f0734feb70a4c250d75bbf0b8a'
    if filename == 'v5-05-52-06_ak.zip':
        return 'c54372a336cd0df69cf7d62aca099c9d8ba4360724c8c07e23aeab94e4c40443'
    if filename == 'v5-04-17-05_nm-cibola.csv':
        return '2e568d0acef546341b26b744cb082dfe126ec2ce516230d1dcd4037a8ebdf4c0'
    if filename == 'v5-05-32-04_tn-williamson.csv':
        return '7e1cf7e9f9d2fcffdda9c1132bcbcd39e8a5a09af791d9299b2bc549ece1b795'
    if filename == 'v5-05-12-01_mi-wayne.zip':
        return '09985058b625d14d87750fb0f727d86df6cda83fda01b56fa69208164fec0a7d'
    if filename == 'v5-15-16-01_nj-cumberland.csv':
        return '00f41061eaa452772c22a120a963d5225c6aecb8a29cf79d2ca20f53bebcf029'
    if filename == 'v5-05-12-01_ga-walker.csv':
        return '7ee53c6519b792f693fd6e5dd25ddfaaf0e4ce8d32cfe92d6bec707b05cd6932'
    if filename == 'v5-05-32-04_az-maricopa.zip':
        return 'ad76cfef492baeb6c8911e1f8fc4e9fba54d4a4f4b469a1270831ca04d3f7d42'

    assert False, 'No expected sanitized hash for %s' % filename
