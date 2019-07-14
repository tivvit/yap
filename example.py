import yap

p = yap.Pipeline(settings={
    "state": {
        "type": "json"
    }
})
p.add(yap.Block("a", "ls"))
bl = yap.DictBlock({"name": "B", "exec": "ls -la"})
p.add(bl)
print(p)
