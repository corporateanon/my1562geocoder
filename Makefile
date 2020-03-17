DATA_FILE_VERSION = v0.0.1

DOWNLOAD_URL = https://github.com/my1562/normalizer/releases/download/$(DATA_FILE_VERSION)/geocoder-data.gob
LOCAL_DIR = data
LOCAL_PATH = $(LOCAL_DIR)/geocoder-data.gob

default:
	mkdir -p $(LOCAL_DIR)
	rm -f $(LOCAL_PATH)
	wget $(DOWNLOAD_URL) -O $(LOCAL_PATH)
	pkger
