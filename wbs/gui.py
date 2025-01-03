import cv2
import numpy as np
import gi
import sys
import datetime

gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, Gdk, GLib, GdkPixbuf

class WBSProject:
    def __init__(self, name, folder, date_created, date_changed):
        self.name = name if name else "Default Project"
        self.folder = folder if folder else "/tmp/default_project"
        self.date_created = date_created if date_created else datetime.now()
        self.date_changed = date_changed if date_changed else datetime.now()


class WBSState:

    def __init__(self):
        self.settings = self._loadSettings()
        self.projects = self._loadProjects()
        self.current_project

    def _loadProjects(self):
        pass

    def _loadCurrentProject(self):
        pass

    def _loadSettings(self):
        pass



class WBS(Gtk.Window):
    
    project_title = "DefaultTODO"
    image_count = 0

    def __init__(self):

        super().__init__(title="Webcam Feed with GTK")
        # TODO start maximized
        self.set_default_size(640, 480)
        
        box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing = 2)
        self.add(box)

        menu_bar = self.__build_main_menu()
        box.pack_start(menu_bar, expand = False, fill = False, padding = 1)
        
        button = Gtk.Button(label="Capture image")
        button.connect("clicked", self.on_button_clicked)
        box.pack_start(button, expand=False, fill=False, padding=1)

        self.viewport = Viewport('/dev/video3')

        box.pack_start(self.viewport, expand=True, fill=True, padding=6)
        
        self.connect("destroy", self.on_destroy)        
        self.show_all()

    def on_button_clicked(self, widget):
        print("prenteding to capture image ... ")
        self.viewport.save_frame(self.image_count, self.project_title)
        self.image_count += 1
        print("finished pretending ... ")

    def on_destroy(self, widget):
        Gtk.main_quit()

    def run(): # static method
        # Create the GTK application and the window
        window = WBS()

        # Start the GTK main loop
        Gtk.main()

    def __build_main_menu(self):

        accelgroup = Gtk.AccelGroup()
        self.add_accel_group(accelgroup)

        menu_bar = Gtk.MenuBar()
        menu_bar.set_hexpand(True)
                
        menuitem_project = Gtk.MenuItem(label="Project")
        menuitem_project.add_accelerator("activate",
                            accelgroup,
                            Gdk.keyval_from_name("p"),
                            Gdk.ModifierType.MOD1_MASK,
                            Gtk.AccelFlags.VISIBLE)
        menu_bar.append(menuitem_project)

        submenu_project = Gtk.Menu()
        menuitem_project.set_submenu(submenu_project)
        
        # New project
        project_new = Gtk.MenuItem(label="New project")
        project_new.connect('activate', self._on_project_new)
        project_new.add_accelerator("activate", 
                            accelgroup,
                            Gdk.keyval_from_name("n"),
                            Gdk.ModifierType.CONTROL_MASK,
                            Gtk.AccelFlags.VISIBLE)
        submenu_project.append(project_new)

        # Open project
        project_open = Gtk.MenuItem(label="Open existing project")
        project_open.connect('activate', self._on_project_open)
        project_open.add_accelerator("activate", 
                            accelgroup,
                            Gdk.keyval_from_name("o"),
                            Gdk.ModifierType.CONTROL_MASK,
                            Gtk.AccelFlags.VISIBLE)
        submenu_project.append(project_open)

        
        # Close project
        project_close = Gtk.MenuItem(label="Close project")
        project_close.connect('activate', self._on_project_close)
        submenu_project.append(project_close)

        # Close project
        settings_edit = Gtk.MenuItem(label="Global Settings")
        settings_edit.connect('activate', self._on_settings_edit)
        submenu_project.append(settings_edit)


        # Close app
        quit_app = Gtk.MenuItem(label="Quit")
        quit_app.connect('activate', self._on_app_quit)
        submenu_project.append(quit_app)
        return menu_bar
    
    def _on_app_quit(self, widget):
        # TODO save project state
        # TODO save settings
        # TODO clean up temp files etc.
        Gtk.main_quit()

    def _on_project_open(self, widget):
        print("On project open")

    def _on_project_close(self, widget):
        print("On project close")

    def _on_project_new(self, widget):
        print("On project new")

    def _on_settings_edit(self, widget):
        print("On settings edit")
        dialog = SettingsDialog(self)
        response = dialog.run()
        dialog.destroy()

class Viewport(Gtk.DrawingArea):
    
    frame = None
    do_draw = False

    def __init__(self, camera_index):

        super().__init__()

        self.connect("draw", self.on_draw)
        # TODO change tick frequency
        self.add_tick_callback(self.on_draw)

        # Open the webcam
        self.webcam = cv2.VideoCapture()  # 0 is the default webcam
        self.webcam.open('/dev/v4l/by-id/usb-046d_Logitech_BRIO_50316219-video-index0')

        if not self.webcam.isOpened():
            print("Error: Could not open webcam.", file=sys.stderr)
            # exit()
        else:
            self.do_draw = True;

    def on_draw(self, widget, cr):
        if (self.do_draw == False):
            return
        
        self.update_frame()
        if self.frame is not None:
            # Create a GdkPixbuf from the current frame
            pixbuf = GdkPixbuf.Pixbuf.new_from_data(
                self.frame.tobytes(),
                GdkPixbuf.Colorspace.RGB,
                False,
                8,
                self.frame.shape[1],
                self.frame.shape[0],
                self.frame.shape[2] * self.frame.shape[1],
            )
            # Draw the image on the drawing area
            Gdk.cairo_set_source_pixbuf(cr, pixbuf, 0, 0)
            cr.paint()
        else: 
            print("video frame is None")

    def update_frame(self):
        ret, frame = self.webcam.read()
        if ret:
            self.frame = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
            self.queue_draw()
        return True  # Return True to keep the timeout active

    def destroy(self):
        self.webcam.release()  # Release the webcam when the window is closed

    def save_frame(self, image_count, image_base_name):
        success = False
        ret, frame = self.webcam.read()
        if not ret:
            print("Failed to read page scan", file=sys.stderr)
        else: 
            cv2.imwrite(f"{image_base_name}_{image_count}.png", frame)
        

class SettingsDialog(Gtk.Dialog):
    def __init__(self, parent):
        super().__init__(title="My Dialog", transient_for=parent, flags=0)
        self.add_buttons(
            Gtk.STOCK_CANCEL, Gtk.ResponseType.CANCEL, Gtk.STOCK_OK, Gtk.ResponseType.OK
        )
        self.set_default_size(150, 100)
        label = Gtk.Label(label="This is a dialog to display additional information")
        box = self.get_content_area()
        box.add(label)
        self.show_all()