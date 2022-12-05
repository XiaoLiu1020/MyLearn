[Mysql日期数据类型，时间类型使用总结](https://www.cnblogs.com/yhtboke/p/5629152.html)

# 概述-日期类型

| 类型        | 范围                          | 格式                    | 大小          |
| --------- | --------------------------- | --------------------- | ----------- |
| DATE      | 2018-10-24                  | YYYY-MM-DD            | 日期值         |
| TIME      | 23:59:59                    | HH\:MM\:SS            | 时间值         |
| YEAR      | 2019                        | YYYY                  | 年份          |
| DATETIME  | 2019-10-25 23:59:59         | YYYY-MM-DD HH\:MM\:SS | 日期和时间       |
| TIMESTAMP | 2019-10-25 23:59:59/1235642 | YYYY-MM-DD HH\:MM\:SS | 混合日期和时间，时间戳 |

\#　建表默认时间设置

    create table new_table(
        'id1' timestamp NOT NULL default CURRENT_TIMESTAMP,
        'id2' datetime NOT NULL default CURRENT_DATETIME;
        )

类似的还是`CURRENT_TIME,CURRENT_YEAR,CURRENT_DATE`

# 日期函数

获取当前日期＋时间（date + time),函数：now()

    select now() 返回datetime类型数据

还有`localtime(), localtimestamp()...`

# Mysql时间戳(Timestamp)函数

current\_timestamp()

**Unix时间戳,日期**转换函数：

unix\_timestamp(),

unix\_timestamp(date),

from\_unixtime(unix\_timestamp)

## 日期date函数, time 函数

获取当前：curdate()

获取当前时间：curtime()

获取当前UTC日期时间函数：utc\_date(), utc\_time(), utc\_timestamp()

**我国本地时间=UTC时间+8小时**

## 普通选取与Extract(选取)函数

普通

    set @dt = '2008-09-10 07:15:30.123456';

    select date(@dt);        -- 2008-09-10 
    select time(@dt);        -- 07:15:30.123456 
    select year(@dt);        -- 2008 
    select quarter(@dt);     -- 3 
    select month(@dt);       -- 9 
    select week(@dt);        -- 36 
    select day(@dt);         -- 10 
    select hour(@dt);        -- 7 
    select minute(@dt);      -- 15 
    select second(@dt);      -- 30 
    select microsecond(@dt); -- 123456

**Extract()函数**

    set @dt = '2008-09-10 07:15:30.123456';

    select extract(year                from @dt); -- 2008 
    select extract(quarter             from @dt); -- 3 
    select extract(month               from @dt); -- 9 
    select extract(week                from @dt); -- 36 
    select extract(day                 from @dt); -- 10 
    select extract(hour                from @dt); -- 7 
    select extract(minute              from @dt); -- 15 
    select extract(second              from @dt); -- 30 
    select extract(microsecond         from @dt); -- 123456

    select extract(year_month          from @dt); -- 200809 
    select extract(day_hour            from @dt); -- 1007 
    select extract(day_minute          from @dt); -- 100715 
    select extract(day_second          from @dt); -- 10071530 
    select extract(day_microsecond     from @dt); -- 10071530123456 
    select extract(hour_minute         from @dt); --    715 
    select extract(hour_second         from @dt); --    71530 
    select extract(hour_microsecond    from @dt); --    71530123456 
    select extract(minute_second       from @dt); --      1530 
    select extract(minute_microsecond from @dt); --      1530123456 
    select extract(second_microsecond from @dt); --        30123456

    MySQL Extract() 函数除了没有date(),time() 的功能外，其他功能一应具全。并且还具有选取‘day_microsecond' 等功能。注意这里不是只选取 day 和 microsecond，而是从日期的 day 部分一直选取到 microsecond 部分。够强悍的吧！

## dayof函数: dayofweek(), dayofmonth(),dayofyear()

返回日期参数在一周一月一年里的第几天

    set @dt = '2008-08-08';

    select dayofweek(@dt);   -- 6 
    select dayofmonth(@dt); -- 8 
    select dayofyear(@dt);   -- 221

