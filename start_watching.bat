@echo off

pip install -r requirements.txt
pip3 install -r requirements.txt
cd python

py test_watcher.py
python test_watcher.py
python3 test_watcher.py
echo camera started
pause