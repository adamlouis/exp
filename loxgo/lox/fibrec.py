import time

def fib(n):
    if n <= 1:
        return n
    return fib(n - 2) + fib(n - 1)

start = time.time()
for i in range(0,30):
    print(fib(i))
end = time.time()

print("TIME")
print((end - start))