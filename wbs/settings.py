import configparser
import sys
import os
import wbs.xdg_manager

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
    config_file_path = os.path.join( wbs.xdg_manager.xdg_config_home(app_name) + 'config.ini')
    print("Config file path: " + config_file_path)
    with open(config_file_path, 'w') as config_file:
        config.write(config_file)

def create_blank_settings_file_template(app_name: str):
    print(wbs.xdg_manager.xdg_config_home(app_name))
    config = configparser.ConfigParser()
    config["WSB"] = {'webcam': " ", "projects_folder": " "}
    # config["WBS"]["webcam"] = ""
    # config["WBS"]["projects_folder"] = ""
    config_file_path = os.path.join(wbs.xdg_manager.xdg_config_home(app_name), '.config.ini')
    with open(config_file_path, "w") as config_file_handle:
        config.write(config_file_handle)
    return config_file_path