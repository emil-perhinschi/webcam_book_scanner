import cv2
import numpy as np
import gi

gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, Gdk, GLib, GdkPixbuf

class WBS(Gtk.Window):


    def __init__(self):

        self.frame = None
        super().__init__(title="Webcam Feed with GTK")
        # TODO start maximized
        self.set_default_size(640, 480)
        
        # TODO change tick frequency
        self.add_tick_callback(self.on_draw)

        # box = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing = 2)
        # self.add(box)

        # button = Gtk.Button(label="Capture image")
        # button.connect("clicked", self.on_button_clicked)
        # box.pack_end(button, expand=False, fill=False, padding=1)

        self.drawing_area = Gtk.DrawingArea()
        self.drawing_area.connect("draw", self.on_draw)
        # box.pack_end(self.drawing_area, expand=False, fill=False, padding=1)
        self.add(self.drawing_area)
        
        # Open the webcam
        self.webcam = cv2.VideoCapture(0)  # 0 is the default webcam
        if not self.webcam.isOpened():
            print("Error: Could not open webcam.")
            exit()

        self.connect("destroy", self.on_destroy)        
        self.show_all()

    def on_button_clicked(self, widget):
        print("prenteding to capture image ... ")

        print("finished pretending ... ")

    def on_draw(self, widget, cr):
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
            self.drawing_area.queue_draw()
        return True  # Return True to keep the timeout active

    def on_destroy(self, widget):
        self.webcam.release()  # Release the webcam when the window is closed
        Gtk.main_quit()

    def run(): # static method
        # Create the GTK application and the window
        window = WBS()

        # Start the GTK main loop
        Gtk.main()


