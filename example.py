import yap

p = yap.Pipeline()

print(repr(p))

p.add(yap.Block("a", "ls"))

print(p)
print(repr(p))
