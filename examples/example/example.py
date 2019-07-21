from context import yap

p = yap.Pipeline(settings={
    "state": {
        "type": "json"
    }
})
p.add(yap.Block("a", "ls", out=["files.txt", "files1.txt"], in_files=["in.txt"]))
bl = yap.DictBlock({"name": "B", "exec": "ls -la", "deps": ["a"], "in_files": ["in.txt"]})
p.add(bl)
print(p)
