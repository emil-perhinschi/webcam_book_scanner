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
    form_text_entries = {}
    app_name = ""


    def __init__(self, app_name):
        Gtk.Window.__init__(self, title="Edit settings")
        # self.connect("destroy-event", lambda x: Gtk.main_quit())
        self.app_name = app_name

        self.set_default_size(600,400)
        self.connect("destroy-event", self._quit_settings_dialog)

        self.form_container = Gtk.Box(spacing=6, orientation=Gtk.Orientation.VERTICAL)
        self.form_container.set_margin_top(20)
        self.form_container.set_margin_left(10)
        self.add(self.form_container)
        
        grid = Gtk.Grid()
        self.form_container.add(grid)


        self.config = wbs.settings.read_settings(self.app_name)
        
        # print(type(self.config).__name__) ConfigParser

        grid_row = 0
        for section in self.config.sections():
            grid.attach(Gtk.Label(""), 0, grid_row, 1, 1) # empty row
            grid_row += 1
            section_label = Gtk.Label(section)
            section_label.set_xalign(0.0)
            grid.attach(section_label, 0, grid_row, 2, 1) # section label row
            grid_row += 1
            grid.attach(Gtk.Label(""), 0, grid_row, 1, 1) # empty row
            grid_row += 1
            print("Section " + section)
            if (not section in self.form_text_entries):
                self.form_text_entries[section] = {}

            for item in self.config[section]:
                # print(item + " = " + self.config[section][item])
                item_label = Gtk.Label(item)
                item_label.set_size_request(80, 20)
                item_label.set_xalign(0.0)
                grid.attach(item_label, 0, grid_row, 1, 1)

                self.form_text_entries[section][item] = Gtk.Entry()
                self.form_text_entries[section][item].set_text(self.config[section][item])
                self.form_text_entries[section][item].set_size_request(500, 20)
                grid.attach(self.form_text_entries[section][item], 1, grid_row, 1, 1)
                grid_row += 1
        
        button_box = Gtk.Box(spacing=6, orientation=Gtk.Orientation.HORIZONTAL)
        self.form_container.add(button_box)

        save_button = Gtk.Button.new_with_mnemonic("_Save")
        save_button.connect("clicked", self._save_settings)
        button_box.add(save_button)
        self.show_all()

    def _save_settings(self, button):

        print("_save_settings called")

        form_data = {}

        for section in self.form_text_entries:
            if (not section in form_data):
                form_data[section] = {}
            
            for item in self.form_text_entries[section]:
                form_data[section][item] = self.form_text_entries[section][item].get_text()
                print(section + ": " + item + " => " + self.form_text_entries[section][item].get_text())
        
        wbs.settings.save_settings(self.app_name, form_data)

    def _quit_settings_dialog(self):
        Gtk.main_quit()

