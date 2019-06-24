import yaml
from yaml import YAMLObject


class Block(YAMLObject):
    # yaml_tag = "/"
    # yaml_loader = yaml.SafeLoader

    # @classmethod
    # def to_yaml(cls, dumper, data):
    #     # ...
    #     return "a"

    def __init__(self, name: str, exec: str, deps=None, out=None):
        if deps is None:
            deps = []
        if out is None:
            out = []
        self.name = name
        self.exec = exec
        self.deps = deps
        self.out = out

    # def __repr__(self) -> str:
    #     return yaml.dump({
    #         "name": self.name,
    #         "exec": self.exec,
    #         "deps": self.deps,
    #         "out": self.out,
    #     }, default_flow_style=False)


class DictBlock(Block):
    def __init__(self, name: str, params: dict):
        exec = params.get("exec", "")
        if not exec:
            raise Exception("Missing exec param")
        super(DictBlock, self).__init__(name, exec, **params)


class Pipeline(object):
    def __init__(self):
        self.pipeline = {}

    def __repr__(self) -> str:
        # o = {k: repr(v) for k, v in self.pipeline.items()}
        return yaml.dump(self.pipeline, default_flow_style=False, tags=False)

    def load_from_file(self, fn: str):
        for k, v in yaml.load(open(fn, "r")):
            self.pipeline[k] = DictBlock(k, v)

    def add(self, block: Block):
        self.pipeline[block.name] = block
