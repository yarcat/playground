```bash
sysctl -n machdep.cpu.brand_string
Apple M4

# With boundary checks:
go run .
                 min(ms)   50(ms)  max(ms)
     simpleLoop    48.37    50.75    54.89
     indepSums2    37.20    37.23    37.76
     indepSums4    24.79    24.81    25.56
     indepSums8    15.33    15.37    16.65
          paral     7.17     7.21     8.99

# Without boundary checks:
go run -gcflags="-B" .
                 min(ms)   50(ms)  max(ms)
     simpleLoop    50.44    51.01    58.26
     indepSums2    31.01    31.06    32.02
     indepSums4    18.25    18.33    19.48
     indepSums8    14.35    14.40    15.59
          paral     7.02     7.11     8.91
```
