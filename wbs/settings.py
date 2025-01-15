import configparser as cp
import sys

def init_settings_file(app_name: str):
    # config = cp.
    return 0

def settings_file(app_name: str):
    pass

def read_settings(app_name):
    config = cp.ConfigParser()

    # config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')
    if (not config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')):
        print("Failed to read config file.", file=sys.stderr)
        return False
    else:     
        return config