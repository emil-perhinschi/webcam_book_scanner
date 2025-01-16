import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, Gdk, GLib, GdkPixbuf

# https://python-gtk-3-tutorial.readthedocs.io/en/latest/dialogs.html#dialogs

# Subwindow
# https://stackoverflow.com/questions/30991885/close-sub-window-without-closing-main-window-pygtk-in-python

import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk
import wbs.settings

class SettingsDialog(Gtk.Window):
    form_data = {}


    def __init__(self, app_name):
        Gtk.Window.__init__(self, title="Edit settings")
        # self.connect("destroy-event", lambda x: Gtk.main_quit())
        self.set_default_size(600,400)
        self.connect("destroy-event", self._quit_settings_dialog)

        self.form_container = Gtk.Box(spacing=6, orientation=Gtk.Orientation.VERTICAL)
        self.form_container.set_margin_top(20)
        self.form_container.set_margin_left(10)
        self.form_container.add(Gtk.Label("Settings"))
        self.add(self.form_container)
        
        grid = Gtk.Grid()
        self.form_container.add(grid)

        print("appname is " + app_name)
        config = wbs.settings.read_settings(app_name)
        # print(type(self.config).__name__) ConfigParser

        grid_row = 0
        for section in config.sections():
        #     section_box = Gtk.Box(spacing=6, orientation=Gtk.Orientation.VERTICAL)
        #     self.form_container.add(section_box)
        #     section_box.add(Gtk.Label(section))
            # empty row
            grid.attach(Gtk.Label(""), 0, grid_row, 1, 1)
            grid_row += 1
            grid.attach(Gtk.Label(section), 0, grid_row, 1, 1)
            grid_row += 1
            # empty row
            grid.attach(Gtk.Label(""), 0, grid_row, 1, 1)
            grid_row += 1
            print("Section " + section)
            for item in config[section]:
                # print(item + " = " + config[section][item])
                item_label = Gtk.Label(item)
                item_label.set_size_request(80, 20)
                item_label.set_xalign(0.0)
                grid.attach(item_label, 0, grid_row, 1, 1)
                
                if (not section in self.form_data):
                    self.form_data[section] = {}

                self.form_data[section][item] = Gtk.Entry()
                self.form_data[section][item].set_text(config[section][item])
                self.form_data[section][item].set_size_request(500, 20)
                grid.attach(self.form_data[section][item], 1, grid_row, 1, 1)
                grid_row += 1
        
        self.show_all()

    def _quit_settings_dialog(self):
        Gtk.main_quit()

