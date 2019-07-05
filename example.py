import yap

p = yap.Pipeline()
p.add(yap.Block("a", "ls"))
bl = yap.DictBlock({"name": "B", "exec": "ls -la"})
p.add(bl)
print(p)
