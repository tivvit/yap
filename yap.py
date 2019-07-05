import copy

import yaml


class Block(object):
    def __init__(self, name: str, exec: str, check=None, desc=None, deps=None, out=None):
        if deps is None:
            deps = []
        if out is None:
            out = []
        self.name = name
        self.desc = desc
        self.check = check
        self.exec = exec
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
    def __init__(self):
        self.pipeline = {}

    def __repr__(self) -> str:
        return yaml.safe_dump(self.pipeline, default_flow_style=False)

    def load_from_file(self, fn: str):
        for k, v in yaml.load(open(fn, "r")):
            self.pipeline[k] = DictBlock(v)

    def add(self, block: Block):
        self.pipeline[block.name] = block


yaml.add_multi_representer(Block, yaml.dumper.Representer.represent_dict)
yaml.SafeDumper.add_multi_representer(Block, yaml.SafeDumper.represent_dict)
