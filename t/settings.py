from wbs.settings import *
from TAP.Simple import *

plan(1)

ok('ConfigParser' == read_settings('test'), 'config object is ConfigParser')