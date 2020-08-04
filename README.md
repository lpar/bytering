
# bytering

Example Go implementation of a circular byte buffer, for scanning
large slow data streams to find particular byte sequences, without reading all
the data into RAM.

If you can read the data into memory, if the search string is large, or if reading data is relatively fast compared to comparison operations, then you 
will likely do better by implementing [Boyer-Moore string search](https://en.wikipedia.org/wiki/Boyer–Moore_string-search_algorithm).
