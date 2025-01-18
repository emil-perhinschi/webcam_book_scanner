from wbs.settings import *
import TAP.Simple as t
import os

t.plan(1)

config = read_settings('test')
t.ok('ConfigParser' == type(config).__name__, 'config object is ConfigParser')
test_config_file_path = wbs.settings.create_blank_settings_file_template("test_one")
expected_file_path = os.path.join(os.getenv('HOME'), '.config/test_one/.config.ini')
t.ok(test_config_file_path == expected_file_path, 'config file path as expected')
