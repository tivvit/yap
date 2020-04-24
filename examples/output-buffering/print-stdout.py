# if __name__ == '__main__':
#     for i in range(5):
#         print("foo")
import time


def main():
    s = time.time()
    for f in ["a", "b"]:
        print(f"Processing {f}")
    print(f"Total time {time.time() - s:.2f}s")


# if __name__ == '__main__':
main()
