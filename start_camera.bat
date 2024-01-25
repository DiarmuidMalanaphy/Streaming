@echo off

pip install -r requirements.txt
pip3 install -r requirements.txt
cd python

py test_camera.py
python test_camera.py
python3 test_camera.py
echo camera started
pause