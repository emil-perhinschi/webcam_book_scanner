import os, sys, shutil

from gi.repository import GLib


"""
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

def check_xdg_folders(app_name: str):
    return xdg_data_home(app_name) and xdg_cache_home(app_name) and xdg_config_home(app_name)

def remove_xdg_folders(app_name: str):
    home = get_home_path()
    data_home = os.path.join(home, '.local/share', app_name)
    config_home = os.path.join(home, '.config', app_name)
    cache_home = os.path.join(home, '.cache', app_name)
    xdg_folders = [data_home, config_home, cache_home]
    success = True
    for f in xdg_folders:
        if (os.path.isdir(f)):
            shutil.rmtree(f)
        else:
            print("Could not find xdg folder: " + f, file=sys.stderr)
            success = False

    return success


def xdg_data_home(app_name: str):
    home = get_home_path()
    data_home = os.path.join(home, '.local/share', app_name)

    if (not os.path.isdir(data_home)):
        print(data_home + " not found, creating ...")
        os.makedirs(data_home)

    return data_home

def xdg_config_home(app_name: str):
    home = get_home_path()
    config_home = os.path.join(home, '.config', app_name)

    if not os.path.isdir(config_home):
        print(config_home + " not found, creating ...")
        os.makedirs(config_home)

    return config_home

def xdg_cache_home(app_name: str):
    home = get_home_path()
    cache_home = os.path.join(home, '.cache', app_name)

    if not os.path.isdir(cache_home):
        print(cache_home + " not found, creating ...")
        os.makedirs(cache_home)

    return cache_home




