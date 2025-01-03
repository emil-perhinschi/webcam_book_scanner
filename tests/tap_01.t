#!/usr/bin/env python

from TAP.Simple import *
from wbs.xdg_manager import *

plan(3) 
# TODO fix tapsimple to work without a plan initialized
# TODO add "done_testing" to tapsimple

ok("/home/emilper" == get_home_path(), "home path found")

# validate_xdg_data_home("test")

ok(1, "one is one")
