import gi

gi.require_version("Gtk", "3.0")
from gi.repository import Gtk

class MyWindow(Gtk.Window):

    def __init__(self):

        super().__init__(title="Webcam book scanner")
        self.button = Gtk.Button(label="Capture image")
        self.button.connect("clicked", self.on_button_clicked)
        self.add(self.button)


    def on_button_clicked(self, widget):

        print("prenteding to capture image")



win = MyWindow()
win.connect("destroy", Gtk.main_quit)
win.show_all()

Gtk.main()

