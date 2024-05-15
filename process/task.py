import os
import random
import time


def main():
    print("My PID is:", os.getpid())
    time.sleep(random.randint(1, 5))
    exit(random.randint(0, 2))


if __name__ == "__main__":
    main()
