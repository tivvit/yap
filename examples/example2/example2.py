import yap

p = yap.Pipeline()

print(p)

p.add(yap.Block("a", "ls"))

print(p)

bl = yap.DictBlock({"name": "B", "exec": "ls"})
print(bl)

p.add(bl)

print(p)

# yap.Pipeline().load_from_file("a.yml")