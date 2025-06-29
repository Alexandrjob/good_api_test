CREATE DATABASE IF NOT EXISTS testdb;

CREATE TABLE testdb.goods_log (
    Id UInt64,
    ProjectId UInt64,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    EventTime DateTime
) ENGINE = MergeTree() 
ORDER BY EventTime;
