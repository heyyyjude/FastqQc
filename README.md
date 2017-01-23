Fastq stat tool based on CLI
============================

### Usage

FastQc -fq|-gz either a fastq.gz|fastq

####
FastQc file is a binary file for linux. you can just download it and use it.

the quality bar plot meaning  
5 : ----|----| ( mean of 1,2,3,4,5 position quality score )  
10 : ----|----|-- ( mean of 6,7,8,9,10 position quality score )  

### OUTPUT

----Nucleotide stat----  
Total basepair : 1.520501861e+09  
GC ratio : 38.2 %  
----Reads stat----  
Total reads: 5141491  
Max read length: 301  
Min read length: 35  
Median read length: 295.73  
Mean read length: 295.73  
N25: 300  
N50: 301  
N75: 301  
----Quality stat----  
Mean of each postion quality score  
34,34,34,34,34,37,37,37,37,37....up to end of the quality score  
Median Quality score: 29.59  
Mean Quality score: 29.59  
N25: 23.37  
N50: 32.04  
N75: 36.45  
Q30: 65.63 %  
Q20: 72.35 %  
---quality plot---  
5:----|----|----|----|----|----|----  
10:----|----|----|----|----|----|----|--  
15:----|----|----|----|----|----|----|--  
20:----|----|----|----|----|----|----|--  
25:----|----|----|----|----|----|----|--  
30:----|----|----|----|----|----|----|--  
...  
...  
300:----|----|----|-  
