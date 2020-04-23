import sys
import time

if __name__ == '__main__':
    for i in range(5):
        print("foo")
        time.sleep(1)
        print("bar", file=sys.stderr)
