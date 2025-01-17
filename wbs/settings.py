import configparser
import sys

def init_settings_file(app_name: str):
    # config = cp.
    return 0

def settings_file(app_name: str):
    pass

def read_settings(app_name):
    config = configparser.ConfigParser()

    # config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')
    if (not config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')):
        print("Failed to read config file.", file=sys.stderr)
        return False
    else:     
        return config


def save_settings(app_name, form_data):
    
    print(form_data)

    config = configparser.ConfigParser()

    # config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')
    if (not config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')):
        print("Failed to read config file.", file=sys.stderr)
        return False
    
    config = configparser.ConfigParser()

    for section in form_data:
        if (not section in config):
            config[section] = {}

        for item in form_data[section]:
            print(section + ": " + item + " => " + form_data[section][item])
            config[section][item] = form_data[section][item]

    # TODO get the file from the xdg folders 
    # TODO create a template of the config file
    # TODO pick projects folder with a file picker dialog, maybe 
    config_file_path = '/home/emilper/personal/book_scanner/wbs/tests/config.ini'
    with open(config_file_path, 'w') as config_file:
        config.write(config_file)