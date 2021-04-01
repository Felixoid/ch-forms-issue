CREATE TABLE default.test
(
  `Path` String CODEC(ZSTD(3)),
  `Value` Float64 CODEC(Gorilla, LZ4),
  `Time` UInt32 CODEC(DoubleDelta, LZ4),
  `Date` Date CODEC(DoubleDelta, LZ4),
  `Timestamp` UInt32 CODEC(DoubleDelta, LZ4)
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(Date)
ORDER BY (Path, Time)
SETTINGS index_granularity = 256
