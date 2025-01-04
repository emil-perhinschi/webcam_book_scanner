#!/usr/bin/env python

from TAP.Simple import *
from wbs.xdg_manager import *

plan(5) 
# TODO fix tapsimple to work without a plan initialized
# TODO add "done_testing" to tapsimple

eq_ok("/tmp", get_home_path(), "home path found")
eq_ok("/tmp/.local/share/test", validate_xdg_data_home("test"), "data home path found")
eq_ok("/tmp/.config/test", validate_xdg_config_home("test"), "config home path found")
eq_ok("/tmp/.cache/test", validate_xdg_cache_home("test"), "cache home path found")
ok(remove_xdg_folders('test'), 'remove xdg folders')


