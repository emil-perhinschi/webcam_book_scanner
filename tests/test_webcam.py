#!/usr/bin/env python
import unittest
import cv2 as cv
import sys
import os
sys.path.append(os.path.abspath("."))



from webcam import Camera

class TestCamera(unittest.TestCase):

    def test_read_image(self):
        webcam = Camera(0)
        # test camera object instantiates OK
        self.assertEqual(type(webcam).__name__, "Camera", "test camera object instantiated")
        ret, image = webcam.read_image()
        print(type(image))
        print(image)

if __name__ == '__main__':
    unittest.main()