#!/usr/bin/env python

from TAP.Simple import *
from wbs.camera_info import *
plan(1)

ok(list_video_devices()[0] == '/dev/video0', "video0 found")

