@echo off

pip install -r requirements.txt
pip3 install -r requirements.txt
cd ..
cd "Server Code"
cd "Python Client"
 

py test_watcher.py
python test_watcher.py
python3 test_watcher.py
echo camera started
pause