import copy

import yaml

VERSION = 1.0

class Block(object):
    def __init__(self, name: str, exe: str, check=None, desc=None, deps=None, out=None):
        if deps is None:
            deps = []
        if out is None:
            out = []
        self.name = name
        self.desc = desc
        self.check = check
        self.exec = exe
        self.deps = deps
        self.out = out

    def items(self):
        r = copy.deepcopy(self.__dict__)
        del r["name"]
        return r.items()

    def __repr__(self) -> str:
        return yaml.dump(self.__dict__, default_flow_style=False)


class DictBlock(Block):
    def __init__(self, params: dict):
        name = params.get("name", "")
        if not name:
            raise Exception("Missing name param")
        exe = params.get("exec", "")
        if not exe:
            raise Exception("Missing exec param")
        del params["name"]
        del params["exec"]
        super(DictBlock, self).__init__(name, exe, **params)


class Pipeline(object):
    def __init__(self, settings=None):
        if settings is None:
            settings = {}
        self.pipeline = {}
        self.settings = settings

    def __repr__(self) -> str:
        out = {
            "version": VERSION,
            "pipeline": self.pipeline,
            "settings": self.settings
        }
        return yaml.safe_dump(out, default_flow_style=False)

    def load_from_file(self, fn: str):
        for k, v in yaml.load(open(fn, "r")):
            self.pipeline[k] = DictBlock(v)

    def add(self, block: Block):
        self.pipeline[block.name] = block


yaml.add_multi_representer(Block, yaml.dumper.Representer.represent_dict)
yaml.SafeDumper.add_multi_representer(Block, yaml.SafeDumper.represent_dict)
