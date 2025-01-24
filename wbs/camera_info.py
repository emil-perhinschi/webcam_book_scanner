import os
import glob
import re
import subprocess

def list_video_devices():
    result = []
    for device in glob.glob('/dev/v4l/by-id/usb*'):
        can_do_video = test_device_capabilities(device)
        if (can_do_video):
            result.append(device)

    return result

def test_device_capabilities(device: str):
    result = subprocess.run(['v4l2-ctl', '-D', '--device=' + device], capture_output=True).stdout.decode('UTF-8')
    
    regex = re.compile('\tDevice Caps.*\n\t\tVideo Capture')
    if(re.search(regex, result)):
        return True
    else:
        return False



# print(list_video_devices())