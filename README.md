# fio-benchmark

fio-benchmark is a wrapper for [fio](https://github.com/axboe/fio) benchmarks. Fio (Flexible I/O Tester) is a tool for storage performance benchmarking.

## Building and running

### Test and Build

```
make test
make build 
```

### Running

Running the fio-benchmark requires `fio` and the `libaio` development packages to be installed on the host.

```
bin/fio-benchmark <flags>
```

#### Usage

```
bin/fio-benchmark -h
```

#### Flags
| Name            |  Description |
|-----------------|--------------------------------------------------------------------------------------------------|
| --output-file   | redirect fio benchmark result to output file                                                     |
| --render-format | redirect fio benchmark result to output file with rendered format, eg. table, html, markdown, csv|
| --config-file   | fio benchmark config file                                                                        |
| --dryrun        | dry-run (default true)                                                                           |
| --v             | number for the log level verbosity                                                               |

### Config file
```yaml
fio_settings:
  numjobs: # 1 2 4 8 16 32 64 128 256 512 1024 2048
  - 1
  - 2
  - 4
  - 8
  - 16
  - 32
  - 64
  - 128
  - 512
  ioengine: libaio
  direct: true
  verify: true
  bs: # block size 4K, 8K, 16K, 32K, 256K, 512K, 1M, 4M
  - 4K
  - 16K
  - 32K 
  - 256K 
  - 1M
  - 4M
  runtime: 15 # seconds
  iodepth:  # 1, 2, 4, 8, 16, 32, 64, 128
  - 1
  - 4
  - 8
  - 16
  - 32
  rw: # read, write, randread, randwrite, rw, randrw
  - read
  - write
  - randread
  - randwrite
  - rw
  - randrw
  # filename: /dev/vdb # device name or file name, which can be ignore if specify `use_all_disks`
use_all_disks: true # except root disk
workers: 8 # It is recommended to be less than or equal to the number of disks
```

## Output
| filename | numjobs | runtime | ioengine | direct | verify | blocksize | iodepth | rw | read-iops-mean | read-bw-mean(kiB/s) | latency-read-min(us) | latency-read-max(us) | latency-read-mean(us) | read-stddev(us) | write-iops-mean | write-bw-mean(kiB/s) | latency-write-min(us) | latency-write-max(us) | latency-write-mean(us) | latency-write-stddev(us) |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | ---:| ---:| ---:| ---:| ---:| ---:| ---:| ---:| ---:| ---:| ---:| ---:|
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4K | 1 | randread | 1112 | 4448 | 321 | 214401 | 900.86086813 | 2172.317800256 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4K | 1 | randrw | 383.789474 | 1535.157895 | 371 | 9878 | 754.054627583 | 437.628838461 | 396.631579 | 1586.526316 | 829 | 268986 | 1767.501507166 | 4364.034579003 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4K | 8 | randread | 9769.684211 | 39078.736842 | 238 | 59433 | 813.701368859 | 1051.125752575 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4K | 8 | randrw | 1812.736842 | 7250.947368 | 269 | 224942 | 1439.133288492 | 2066.725188392 | 1811.473684 | 7245.894737 | 647 | 233832 | 2964.4548968199997 | 4637.792082346 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4M | 1 | randread | 191.210526 | 783330.368421 | 2867 | 19450 | 5201.692435348999 | 1859.916525032 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4M | 1 | randrw | 67.052632 | 274647.578947 | 3105 | 18326 | 5037.5913099849995 | 1703.040075248 | 71.052632 | 291031.578947 | 5135 | 28093 | 9220.460498602999 | 2244.186007376 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4M | 8 | randread | 583.157895 | 2.388614736842e+06 | 5864 | 67574 | 13681.189661591001 | 2738.0083960859997 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 1 | 10s | libaio | 1 |  | 4M | 8 | randrw | 162.25 | 664757.9 | 3388 | 239363 | 21933.212207744 | 13677.558721186 | 168.05 | 688507.35 | 6924 | 241056 | 26336.419029691 | 17291.395360652 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4K | 1 | randread | 10594.210526 | 42377.157895 | 260 | 21152 | 768.719952224 | 730.526613278 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4K | 1 | randrw | 1956.421053 | 7825.684211 | 286 | 13070 | 1332.274188352 | 1060.184858525 | 1985.473684 | 7941.894737 | 599 | 19227 | 2705.696762644 | 1305.309295732 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4K | 8 | randread | 44452.210526 | 177809.789474 | 283 | 268095 | 1426.7336875370002 | 4108.448719071 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4K | 8 | randrw | 4053.052632 | 16212.210526 | 312 | 41883 | 6971.509941483 | 4883.835857048 | 4089.684211 | 16358.736842 | 567 | 46525 | 8824.19664812 | 4872.5854249 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4M | 1 | randread | 577.971053 | 2.368886626316e+06 | 4165 | 49397 | 13852.187700086999 | 3027.6172978170002 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4M | 1 | randrw | 164.473684 | 676696.842105 | 4633 | 229249 | 21780.814631868 | 12950.23341288 | 165.947368 | 682680.789474 | 9080 | 439031 | 26270.942311796 | 21726.045013228 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4M | 8 | randread | 570.1 | 2.3351583e+06 | 28271 | 363925 | 112130.887111979 | 35518.905668209 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 8 | 10s | libaio | 1 |  | 4M | 8 | randrw | 197.4 | 809305.5 | 5755 | 376474 | 96048.028940733 | 29888.987186832 | 201.3 | 825289.7 | 40001 | 1064964 | 221832.835930549 | 103129.34279955301 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4K | 1 | randread | 47495.052632 | 189994.157895 | 282 | 317402 | 1339.909596079 | 3429.103578518 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4K | 1 | randrw | 3755.931579 | 15043.115789 | 362 | 34836 | 7503.007442663 | 5188.695334104 | 3777.536842 | 15129.536842 | 740 | 42328 | 9423.755584665001 | 5118.7020513709995 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4K | 8 | randread | 59435.842105 | 237760.210526 | 355 | 79744 | 8547.212181693001 | 5576.0523857260005 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4K | 8 | randrw | 8639.528947 | 34581.815789 | 1127 | 336468 | 21886.2323931 | 19443.902509673 | 8613.823684 | 34479.244737 | 2189 | 351591 | 37302.050208586 | 24028.957604481002 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4M | 1 | randread | 575.739474 | 2.3593553e+06 | 12442 | 519580 | 110518.32612023399 | 40456.831552032 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4M | 1 | randrw | 251.606468 | 1.066944513234e+06 | 21014 | 380358 | 94714.87367759299 | 34632.678864924 | 198.856742 | 850408.518894 | 36772 | 1105923 | 229192.196654647 | 110551.67751639799 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4M | 8 | randread | 573.762711 | 2.35014751302e+06 | 23993 | 1871285 | 851647.256101975 | 210761.694246333 | 0 | 0 | 0 | 0 | 0 | 0 |
| /dev/vdb | 64 | 10s | libaio | 1 |  | 4M | 8 | randrw | 243.755091 | 1.034801559985e+06 | 84499 | 2631350 | 1.074421372005806e+06 | 359340.135599382 | 239.719699 | 1.020531372308e+06 | 84170 | 2732242 | 1.184413206235976e+06 | 368635.16110135696 |
