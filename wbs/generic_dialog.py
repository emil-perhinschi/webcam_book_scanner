import os
import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, Gdk, GLib, GdkPixbuf

import wbs.settings

class GenericDialog(Gtk.Window):
    form_text_entries = {}
    app_name = ""

    def __init__(self, app_name, dialog_title):
        if (not dialog_title):
            dialog_title = "Default dialog title"
        
        Gtk.Window.__init__(self, title = dialog_title)
        self.app_name = app_name

        self.set_default_size(600,400)
        self.connect("destroy-event", self._quit_settings_dialog)
        
        self.config = self.load_config_data()

        self.form_container = Gtk.Box(spacing=6, orientation=Gtk.Orientation.VERTICAL)
        self.form_container.set_margin_top(20)
        self.form_container.set_margin_left(10)
        self.add(self.form_container)
        
        grid = Gtk.Grid()
        self.form_container.add(grid)
        settings_file_path = wbs.settings.settings_file(app_name)
        if (not os.path.isfile(settings_file_path)):
            wbs.settings.create_blank_settings_file_template(app_name)

        self.attach_form(grid)

        button_box = Gtk.Box(spacing=6, orientation=Gtk.Orientation.HORIZONTAL)
        self.form_container.add(button_box)

        save_button = Gtk.Button.new_with_mnemonic("_Save")
        save_button.connect("clicked", self._save_form)
        button_box.add(save_button)
        self.show_all()

    def load_config_data(self):
        raise Exception("override me")


    def attach_form(self, grid): 
        raise Exception("override me")

    def _attach_section_label(self, grid: Gtk.Grid, grid_row: int, section: str):
        grid.attach(Gtk.Label(""), 0, grid_row, 1, 1) # empty row
        grid_row += 1
        section_label = Gtk.Label(section)
        section_label.set_xalign(0.0)
        grid.attach(section_label, 0, grid_row, 2, 1) # section label row
        grid_row += 1
        grid.attach(Gtk.Label(""), 0, grid_row, 1, 1) # empty row
        grid_row += 1


    def _attach_text_entry(self, grid: Gtk.Grid, grid_row: int, section: str, item: str): 
        item_label = Gtk.Label(item)
        item_label.set_size_request(80, 20)
        item_label.set_xalign(0.0)
        grid.attach(item_label, 0, grid_row, 1, 1)

        if (not section in self.form_text_entries):
            self.form_text_entries[section] = {}

        self.form_text_entries[section][item] = Gtk.Entry()
        self.form_text_entries[section][item].set_text(self.config[section][item])
        self.form_text_entries[section][item].set_size_request(500, 20)

        grid.attach(self.form_text_entries[section][item], 1, grid_row, 1, 1)


    def _attach_folder_entry(self, grid: Gtk.Grid, grid_row: int, section: str, item: str): 
        self._attach_text_entry(grid, grid_row, section, item)
        
        dialog_button = Gtk.Button("Choose folder")
        dialog_button.connect("clicked", self._pick_projects_folder)
        grid.attach(dialog_button, 2, grid_row, 1, 1)

    def _pick_projects_folder(self, button):
        dialog = Gtk.FileChooserDialog(
            title="Please choose a folder", 
            parent=self, 
            action=Gtk.FileChooserAction.SELECT_FOLDER
        )
        dialog.set_current_folder(GLib.get_home_dir())
        dialog.add_buttons(
            Gtk.STOCK_CANCEL,
            Gtk.ResponseType.CANCEL,
            Gtk.STOCK_OPEN,
            Gtk.ResponseType.OK,
        )

        # self.add_projects_folder_filters(dialog)
        response = dialog.run()
        if response == Gtk.ResponseType.OK:
            print("Open clicked")
            print("File selected: " + dialog.get_filename())
            self.form_text_entries["WBS"]["projects_folder"].set_text(dialog.get_filename())
        elif response == Gtk.ResponseType.CANCEL:
            print("Cancel clicked")

        dialog.destroy()

    def _save_form(self, button):

        form_data = {}

        for section in self.form_text_entries:
            if (not section in form_data):
                form_data[section] = {}
            
            for item in self.form_text_entries[section]:
                form_data[section][item] = self.form_text_entries[section][item].get_text()
                print(section + ": " + item + " => " + self.form_text_entries[section][item].get_text())
    
        self._save_form_data_to_file(form_data)

    def _save_form_data_to_file(self, form_data):
        raise Exception("override me")



    def _quit_settings_dialog(self):
        Gtk.main_quit()

