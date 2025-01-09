import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, Gdk, GLib, GdkPixbuf

# https://python-gtk-3-tutorial.readthedocs.io/en/latest/dialogs.html#dialogs

# Subwindow
# https://stackoverflow.com/questions/30991885/close-sub-window-without-closing-main-window-pygtk-in-python

import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk


class SettingsDialog(Gtk.Window):
    def __init__(self):
        Gtk.Window.__init__(self, title="Edit settings")
        # self.connect("destroy-event", lambda x: Gtk.main_quit())
        self.connect("destroy-event", self._quit_settings_dialog)

        self.add(Gtk.Label("Settings"))
        self.show_all()

    def _quit_settings_dialog(self):
        Gtk.main_quit()