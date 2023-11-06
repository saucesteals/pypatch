import os
print("Hello from inject.py!")
print = lambda *args, **kwargs: os.write(1, b"Boo!\n")