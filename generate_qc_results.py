import glob
import os
import sys
from collections import OrderedDict

__author__ = "jkkim"


def get_out_file_list(path):
    glob_path = os.path.join(path, "*.out")
    return [i for i in glob.glob(glob_path) if os.path.isfile(i)]


def parse_out(outf):
    result = OrderedDict()
    with open(outf, 'rt')as fin:
        for line in fin:
            line = line.strip()

            if line.startswith("Total basepair"):
                result["TotalBasepair"] = line.split(":")[1].strip()

            if line.startswith("Total reads"):
                result["TotalReads"] = line.split(":")[1].strip()

            if line.startswith("GC ratio"):
                result["GCRatio"] = line.split(":")[1].strip().replace(" %",
                                                                       "%")
            if line.startswith("Mean Quality score"):
                result["MeanQualityScore"] = line.split(":")[
                    1].strip().replace(" %", "")

            if line.startswith("Q30"):
                result["Q30"] = line.split(":")[1].strip().replace(" %", "%")

            if line.startswith("Q20"):
                result["Q20"] = line.split(":")[1].strip().replace(" %", "%")

    return result


def show_result(file_name, res_dict):
    '''
    TotalBasepair
    GCRatio
    TotalReads
    MeanQualityScore
    Q30
    Q20
    '''

    file_name = file_name.replace(".fastq.gz.out", "")
    result = list()
    result.append(file_name)
    for k in res_dict:
        print(k)
        result.append(res_dict[k])
    out = ",".join(result)
    print(out)


def main(path):
    dirs = get_out_file_list(path)

    for infile in dirs:
        show_result(infile, parse_out(infile))


if __name__ == '__main__':
    path = sys.argv[1]
    main(path)
