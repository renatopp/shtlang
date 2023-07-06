

def fib(n):
  i, a, b = 0, 0, 1

  if n >= 0:
    yield a
  
  if n >= 1:
    yield b
  
  while i < n:
    a, b = b, a + b
    i += 1
    yield b

for i in fib(10):
  print(i)