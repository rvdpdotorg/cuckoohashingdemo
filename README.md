This is mainly an exercise in [Go](https://golang.org/) programming. I have
started to learn Go and _cuckoo hashing_ seemed a nice exercise to implement
in Go.

I have written a
[blog](https://kirk.rvdp.org/homepage/posts/2019/cuckoo-hashing/)
about how cuckooo hashing works and why it is a useful datastructure.

The program asks for the number of hash bits to use for the hash
table and whether you want to do cuckoo hashing or not. The result
is printed on _stdout_. The first column is the amount of elements
to store in the hash table. The second column is the hash
collision probability.
```
$ ./cuckoohashingdemo
Number of hash bits (hash table size 2^bits): 11
Use Cuckoo hashing (y/n)? y
Number of samples: [2] [66] [130] [194] [258] [322] [386] [450] [514] [578] [642] [706] [770] [834] [898] [962] [1026] [1090] [1154] [1218] [1282] [1346] [1410] [1474] [1538] [1602] [1666] [1730] [1743] [1756] [1769] [1782] [1795] [1808] [1821] [1834] [1847] [1860] [1873] [1886]
2 0.000000
66 0.000000
130 0.000000
194 0.000000
258 0.000000
322 0.000000
386 0.000000
450 0.000000
514 0.000000
578 0.000000
642 0.000000
706 0.000000
770 0.000000
834 0.000000
898 0.000000
962 0.000000
1026 0.000000
1090 0.000000
1154 0.000000
1218 0.000000
1282 0.000000
1346 0.000000
1410 0.000100
1474 0.000600
1538 0.000700
1602 0.008200
1666 0.033000
1730 0.152600
1743 0.199200
1756 0.264400
1769 0.330100
1782 0.424900
1795 0.526500
1808 0.633000
1821 0.733600
1834 0.829000
1847 0.897900
1860 0.951800
1873 0.981500
1886 0.994600
```

This code is not meant for production.
