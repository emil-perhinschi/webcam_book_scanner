import os
import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, GLib
import wbs.settings
from wbs.generic_dialog import *

class SettingsDialog(GenericDialog):
    form_text_entries = {}
    app_name = ""

    def __init__(self, app_name,):
        GenericDialog.__init__(self, app_name=app_name, dialog_title="Modify Project")

    def attach_form(self, grid):
        # FORM
        grid_row = 0
        self._attach_section_label(grid, grid_row, "Project")
        grid_row = 4
        self._attach_text_entry(grid, grid_row, "WBS", "project_name")
        grid_row = 5
        self._attach_folder_entry(grid, grid_row, "WBS", "projects_folder")
        
    def _save_form_data_to_file(self, form_data):
        # TODO 
        pass

