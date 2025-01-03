import os, sys
from gi.repository import GLib


"""
TODO deal with system-wide config, data, cache folders, I guess 
/etc/wbs/
/usr/shared/wbs/
/tmp/wbs/
"""


"""
TODO verify the XDG folders exist, if not create
XDG_DATA_HOME # default $HOME/.local/share, user-specific data 

XDG_CONFIG_HOME # default $HOME/.config, configuration files

XDG_CACHE_HOME # default $HOME/.cache, cache/temporary/disposable files which can be deleted 
"""

def get_home_path():
    # GLib.get_home_dir does not work from within the VSCode terminal
    home = GLib.get_home_dir()

    if (not home):
        home = os.getenv('HOME')

    return home # TODO find other ways to guess the home if these two don't work

def validate_xdg_folders(app_name: str):
    return validate_xdg_data_home(app_name) and validate_xdg_cache_home(app_name) and validate_xdg_config_home(app_name)

def validate_xdg_data_home(app_name: str):
    home = get_home_path()
    data_home = os.path.join(home, '.local/share', app_name)

    if (not os.path.isdir(data_home)):
        print(data_home + " not found, creating ...")
        print(data_home)
        os.makedirs(data_home)


def validate_xdg_config_home(app_name: str):
    home = get_home_path()
    config_home = os.path.join(home, '.config', app_name)

    if not os.path.isdir(config_home):
        print(config_home + " not found, creating ...")
        os.makedirs(config_home)

def validate_xdg_cache_home(app_name: str):
    home = get_home_path()
    cache_home = os.path.join(home, '.cache', app_name)

    if not os.path.isdir(cache_home):
        print(cache_home + " not found, creating ...")
        os.makedirs(cache_home)




