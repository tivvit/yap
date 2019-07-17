from context import yap

p = yap.Pipeline(settings={
    "state": {
        "type": "json"
    }
})
p.add(yap.Block("a", "ls", out=["files.txt", "files1.txt"]))
bl = yap.DictBlock({"name": "B", "exec": "ls -la", "deps": ["a"]})
p.add(bl)
print(p)
