import cv2 as cv
import gi

gi.require_version("Gtk", "3.0")
from gi.repository import Gtk

from webcam import Camera

class WBS(Gtk.Window):

    def run():
        object = WBS()
        object.connect("destroy", Gtk.main_quit)
        object.show_all()
        Gtk.main()

    def __init__(self):
        
        super().__init__(title="Webcam book scanner")

        self.camera = Camera(0)
        self.button = Gtk.Button(label="Capture image")
        self.button.connect("clicked", self.on_button_clicked)
        self.add(self.button)

    def on_button_clicked(self, widget):
        print("prenteding to capture image ... ")
        ret, image = self.camera.read_image()

        print("finished pretending ... ")
        

class CameraViewport (Gtk.Frame):

    def __init__(self, css=None, border_width=1):
        super().__init__()
        self.set_border_width(border_width)
        self.set_size_request(300,300)
        self.vexpand = True
        self.hexpand = True
        self.surface = None
        self.area = Gtk.DrawingArea()
        self.add(self.area)

        self.area.connect("draw", self.on_draw)
        self.area.connect('configure-event', self.on_configure)

    def on_draw(self, area, context):
        if self.surface is not None:
            context.set_source_surface(self.surface, 0.0, 0.0)
            context.paint()
        else:
            print("invalid surface")

        return False
    
    def on_configure(self, area, event, data=None):
        self.redraw()
        return False
    
    # def redraw():
        