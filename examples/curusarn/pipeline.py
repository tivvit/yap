from context import yap
import os

SHELL = "/bin/bash"
DATA = "/mfs/replicated/datasets/nametag-target-mapping"

# global revision root
REV_ROOT = os.path.join(DATA, "/versions")
REV_ROOT_BY_DATE = os.path.join(DATA, "/versions_by_date")
REV_ROOT_BY_NAME = os.path.join(DATA, "/versions_by_name")
# tmp data directory
TMP_DATA = "/tmp/nametag-target-mapping-data"

# where to save data in this run
# nothing to commit or stage => ${REV_ROOT}/"current commit hash"
# uncommited/unstaged data => ${TMP_DATA}
REV_DATA = "{} ./get-revision-data-path {} {}".format(SHELL, REV_ROOT, TMP_DATA)

LAST_COMMIT_DATETIME = SHELL + "git log -1 --format=\"%at\" | xargs -I{} " \
                               "date -d @{} +%Y-%m-%d_%H:%M:%S"
LAST_COMMIT_MSG = SHELL + "git log -1 --pretty=%B | head -n 1 | tr ' ' '_'"

make2graph_svg = os.path.join(REV_DATA, "/make2graph.svg")

# data files definition
transactional_dataset = "/mfs/replicated/datasets/phishing" \
                        "/transactional_dataset/transactional"
fulltextPremises_xml = os.path.join(DATA,
                                    "/firmy-cz-data/2019-01_fulltextPremises.xml")
domainWeights_json = os.path.join(DATA,
                                  "/domain-weights/2019-04_domain-weights.json")

blacklistAdeptsDomains_txt = os.path.join(REV_DATA,
                                          "/blacklistAdeptsDomains.txt")

blacklist_txt = os.path.join(REV_DATA, "/blacklist.txt")
crossHostingBlacklist_txt = "/tmp/nametag-target-mapping-curated-domain-lists" \
                            "/cross-hosting-blacklist.txt"
scannerBypassWhitelist_txt = "/tmp/nametag-target-mapping-curated-domain-lists/scanner-bypass-whitelist.txt"

transactionalDomains_txt = os.path.join(REV_DATA, "/transactionalDomains.txt")
transactionalDomainsUniq_txt = os.path.join(REV_DATA,
                                            "/transactionalDomainsUniq.txt")

importantDomains_txt = os.path.join(REV_DATA, "/importantDomains.txt")

getffqTitle_txt = os.path.join(REV_DATA, "/getffqTitle_txt")
fulltextPremisesTitle_txt = os.path.join(REV_DATA, "/fulltextPremisesTitle.txt")
fulltextPremisesGetffqTitle_txt = os.path.join(REV_DATA,
                                               "/fulltextPremisesGetffqTitle.txt")

suffixFile_json = os.path.join(REV_DATA, "/suffixFile.json")
suffixFile_svg = os.path.join(REV_DATA, "/suffixFile.svg")
tmp_suffixFile_svg = "company-suffixes.svg"

fulltextPremises_json = os.path.join(REV_DATA, "/fulltextPremises.json")
fulltextPremisesFiltered_json = os.path.join(REV_DATA,
                                             "/fulltextPremisesFiltered.json")
fulltextPremisesFilteredWithRoots_json = os.path.join(REV_DATA,
                                                      "/fulltextPremisesFilteredWithRoots.json")
fulltextPremisesMapping_json = os.path.join(REV_DATA,
                                            "/fulltextPremisesMapping.json")
fulltextPremisesReversedMapping_json = os.path.join(REV_DATA,
                                                    "/fulltextPremisesReversedMapping.json")

gen_prefix = ""

# USER FACING TARGETS

p = yap.Pipeline(settings={
    "state": {
        "type": "json"
    }
})

# run whole build
# p.add(yap.Block("run", "true", deps=["fulltextPremisesReversedMapping_json",
#                                      "make2graph.svg"]))

# show current revision path
p.add(yap.Block("rev-path", "echo {}".format(REV_DATA)))
p.add(
    yap.Block("clean", "echo -e Please run this:\nrm -rf {}".format(REV_DATA)))
p.add(yap.Block("list-imports",
                "echo \"PYTHON3\";"
                "cat *.py | grep 'import ' | sort | uniq;"
                "echo \"PYTHON2\";"
                "cat *.py2 | grep 'import ' | sort | uniq"))
# p.add(yap.Block("make2graph", "true", deps=["make2graph.svg"]))

# ACTUAL TARGETS
p.add(yap.Block(gen_prefix + "/tmp/nametag-target-mapping-curated-domain-lists",
                "cd /tmp && git clone "
                "git@gitlab.seznam.net:simon.let/nametag-target-mapping-curated-domain-lists.git",
                out=["/tmp/nametag-target-mapping-curated-domain-lists"]))
p.add(yap.Block(gen_prefix + crossHostingBlacklist_txt.replace('/', '_'),
                "true",
                out=[fulltextPremises_json],
                in_files=["/tmp/nametag-target-mapping-curated-domain-lists"]))
p.add(yap.Block(gen_prefix + scannerBypassWhitelist_txt,
                "true",
                out=[scannerBypassWhitelist_txt.replace('/', '_')],
                in_files=["/tmp/nametag-target-mapping-curated-domain-lists"]))
p.add(yap.Block(gen_prefix + blacklist_txt,
                "bash -c \"cat {} > {}\"".format(crossHostingBlacklist_txt, blacklist_txt),
                in_files=[crossHostingBlacklist_txt],
                out=[blacklist_txt]
                ))
p.add(yap.Block(gen_prefix + importantDomains_txt,
                "python3 getImportantDomains.py --input-file {} > {}".format(
                    domainWeights_json, importantDomains_txt) +
                "cat {} >> {}".format(scannerBypassWhitelist_txt,
                                      importantDomains_txt),
                in_files=["getImportantDomains.py",
                          domainWeights_json,
                          scannerBypassWhitelist_txt],
                out=[importantDomains_txt]))
p.add(yap.Block(gen_prefix + transactionalDomains_txt,
                "python2 getTransactionalDomains.py2 --input-directory {} > {}".format(
                    transactional_dataset,
                    transactionalDomains_txt),
                in_files=["getTransactionalDomains.py2",
                          transactional_dataset],
                out=[transactionalDomains_txt]))
p.add(yap.Block(gen_prefix + transactionalDomainsUniq_txt,
                "cat {} | sort | uniq > {}".format(
                    transactionalDomains_txt,
                    transactionalDomainsUniq_txt),
                in_files=[transactionalDomains_txt],
                out=[transactionalDomainsUniq_txt]))
p.add(yap.Block(gen_prefix + fulltextPremises_json,
                "python3 processFulltextPremises.py --input-file {} > {}".format(
                    fulltextPremises_xml,
                    fulltextPremises_json),
                in_files=["processFulltextPremises.py",
                          fulltextPremises_xml],
                out=[transactionalDomainsUniq_txt]))
c = "python3 processFulltextPremises.py --input-file {} --whitelist {} --create-blacklist-adepts {} --blacklist {} > {}".format(
    fulltextPremises_xml,
    importantDomains_txt,
    blacklistAdeptsDomains_txt,
    blacklist_txt,
    fulltextPremisesFiltered_json)
in_f = ["processFulltextPremises.py",
        fulltextPremises_xml,
        importantDomains_txt,
        blacklist_txt]
out = [fulltextPremisesFiltered_json]
p.add(yap.Block(gen_prefix + fulltextPremisesFiltered_json,
                c,
                in_files=in_f,
                out=out))
p.add(yap.Block(gen_prefix + blacklistAdeptsDomains_txt,
                c,
                in_files=in_f,
                out=out
                ))
p.add(yap.Block(gen_prefix + blacklistAdeptsDomains_txt,
                c,
                in_files=in_f,
                out=out
                ))
p.add(yap.Block(gen_prefix + fulltextPremisesTitle_txt,
                "python3 titleJson2Txt.py --input-file {} > {}".format(
                    fulltextPremises_json,
                    fulltextPremisesTitle_txt),
                in_files=["titleJson2Txt.py",
                          fulltextPremises_json],
                out=[fulltextPremisesTitle_txt]))
p.add(yap.Block(gen_prefix + getffqTitle_txt,
                "scli getffq -f nametag0_if -c -v -s 10000 > {}".format(
                    getffqTitle_txt),
                out=[getffqTitle_txt]))
p.add(yap.Block(gen_prefix + fulltextPremisesGetffqTitle_txt,
                "cat {} {} > {}".format(
                    getffqTitle_txt,
                    fulltextPremisesTitle_txt,
                    fulltextPremisesGetffqTitle_txt,
                ),
                in_files=[getffqTitle_txt,
                          fulltextPremisesTitle_txt],
                out=[fulltextPremisesGetffqTitle_txt]))
p.add(yap.Block(gen_prefix + suffixFile_json,
                "python3 nametag-org-suffix.py -j -l 2 -n 20 -c 30 "
                "--input-file {} > {};".format(
                    fulltextPremisesGetffqTitle_txt,
                    suffixFile_json,
                )
                +
                "python3 nametag-org-suffix.py -g -l 2 -n 20 -c 30 "
                "--input-file {};".format(
                    fulltextPremisesGetffqTitle_txt
                )
                +
                "mv {} {}".format(
                    tmp_suffixFile_svg,
                    suffixFile_svg
                )
                ,
                in_files=["nametag-org-suffix.py",
                          fulltextPremisesGetffqTitle_txt],
                out=[suffixFile_json, suffixFile_svg]))
p.add(yap.Block(gen_prefix + fulltextPremisesFilteredWithRoots_json,
                "python3 getOrgTitleRoot.py -d -s {} -j {} > {}".format(
                    suffixFile_json,
                    fulltextPremisesFiltered_json,
                    fulltextPremisesFilteredWithRoots_json,
                ),
                in_files=["getOrgTitleRoot.py",
                          suffixFile_json,
                          fulltextPremisesFiltered_json],
                out=[fulltextPremisesFilteredWithRoots_json]))
p.add(yap.Block(gen_prefix + fulltextPremisesMapping_json,
                "python3 createMapping.py --input-file {} > {}".format(
                    fulltextPremisesFilteredWithRoots_json,
                    fulltextPremisesMapping_json,
                ),
                in_files=["createMapping.py",
                          fulltextPremisesFilteredWithRoots_json],
                out=[fulltextPremisesMapping_json]))
p.add(yap.Block(gen_prefix + fulltextPremisesReversedMapping_json,
                "python3 reverseMapping.py --input-file {} > {}".format(
                    fulltextPremisesMapping_json,
                    fulltextPremisesReversedMapping_json,
                ),
                in_files=["reverseMapping.py",
                          fulltextPremisesMapping_json],
                out=[fulltextPremisesReversedMapping_json]))

print(p)
