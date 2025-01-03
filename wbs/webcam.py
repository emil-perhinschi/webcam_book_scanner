import cv2 as cv


class Camera:

    def __init__(self, index):

        self.camera = cv.VideoCapture(index)
        self.camera_index = index
        if not self.camera.isOpened():
            print("camera {} could not be opened".format(index))

    def __del__(self):
        self.camera.release()
    
    def read_image(self):

        ret, frame = self.camera.read()
        return ret, cv.cvtColor(frame, cv.COLOR_BGR2RGB)

