@echo off

pip install -r requirements.txt
pip3 install -r requirements.txt
cd ..
cd "Server Code"
cd "Python Client"

py test_camera.py
python test_camera.py
python3 test_camera.py
echo camera started
pause