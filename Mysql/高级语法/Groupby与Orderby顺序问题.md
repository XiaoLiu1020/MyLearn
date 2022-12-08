![表中数据](https://www.linuxidc.com/upload/2019_06/190601200439051.png)
**需求如下：**
查询每个人领取最高奖励，并且*从小到大*排序

    SELECT id, uid, money, datetime FROM reward GROUP by uid ORDER BY DESC;

得到结果如下：
![结果](https://www.linuxidc.com/upload/2019_06/190601200439052.png)

**并没有得到想要结果**

原因：`group by`和`order by`一起使用时，优先使用`group by`分组，并**取出分组后第一条数据**，所以后面`order by`出来都是按照第一条数据排序，**但是第一条数据不一定是最大的数据**

## 解决办法：

### 方法一:子查询

**先排序，再分组，使用子查询**

    SELECT
        r.id,
        r.uid,
        r.money,
        r.datetime
    FROM (SELECT
        id,
        uid,
        money,
        datetime
        FROM reward
        ORDER BY money DESC) r //排序出最大的数据
    GROUP BY r.uid
    ORDER BY r.money DESC;

![方法一结果](https://www.linuxidc.com/upload/2019_06/190601200439053.png)

### 方法二：使用max()

如果不需要取得整条记录，则可以**使用max()**

    SELECT
        id, uid, money, datetime, max(money)
    FROM
        reward 
    GROUP BY 
        uid
    ORDER BY
        MAX(money) DESC;

这只是取出了该uid的最大值，但是并没该最大值的整条数据
