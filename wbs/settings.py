import configparser
import sys
import os
import wbs.xdg_manager

def settings_file(app_name: str):
    config_file_path = os.path.join(wbs.xdg_manager.xdg_config_home(app_name), '.config')
    print(">>>>>>>>>>>>>> " + config_file_path)
    return config_file_path
    

def read_settings(app_name):
    config = configparser.ConfigParser()

    if (not os.path.isfile(settings_file(app_name))):
        create_blank_settings_file_template(app_name)

    # config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')
    if (not config.read(settings_file(app_name))):
        print("Failed to read config file.", file=sys.stderr)
        return False
    else:     
        return config


def save_settings(app_name, form_data):
    
    print(form_data)

    config = configparser.ConfigParser()

    # config.read('/home/emilper/personal/book_scanner/wbs/tests/config.ini')
    if (not config.read(settings_file(app_name))):
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
    config_file_path = settings_file(app_name)
    print("Config file path: " + config_file_path)
    with open(config_file_path, 'w') as config_file:
        config.write(config_file)

def create_blank_settings_file_template(app_name: str):
    print(wbs.xdg_manager.xdg_config_home(app_name))
    config = configparser.ConfigParser()
    config["WBS"] = {'webcam': "", "projects_folder": ""}
    # config["WBS"]["webcam"] = ""
    # config["WBS"]["projects_folder"] = ""
    config_file_path = settings_file(app_name)
    with open(config_file_path, "w") as config_file_handle:
        config.write(config_file_handle)
    print("::::::::::::::::: " + config_file_path)
    return config_file_path