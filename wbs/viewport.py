import cv2
import sys
import numpy as np

import gi
gi.require_version("Gtk", "3.0")
from gi.repository import Gtk, Gdk, GdkPixbuf

class Viewport(Gtk.DrawingArea):
    
    frame = None
    do_draw = False

    def __init__(self):
        super().__init__()



    def connect_to_camera(self, camera_device):
        self.connect("draw", self.on_draw)
        # TODO change tick frequency
        self.add_tick_callback(self.on_draw)

        # Open the webcam
        self.webcam = cv2.VideoCapture()  # 0 is the default webcam
        # self.webcam.open('/dev/v4l/by-id/usb-046d_HD_Pro_Webcam_C920_DC1A8EEF-video-index0')
        camera_device = '/dev/v4l/by-id/usb-046d_Logitech_BRIO_50316219-video-index0'
        self.webcam.open(camera_device)

        if not self.webcam.isOpened():
            print("Error: Could not open webcam.", file=sys.stderr)
        else:
            self.do_draw = True

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
