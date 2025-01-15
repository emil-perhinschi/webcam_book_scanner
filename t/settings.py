from wbs.settings import *
from TAP.Simple import *

plan(1)

config = read_settings('test')
ok('ConfigParser' == type(config).__name__, 'config object is ConfigParser')