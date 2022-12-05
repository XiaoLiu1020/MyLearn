# SUBSTRING

函数从**特定位置开始的字符串返回一个给定长度的子字符串**

```{SQL}
Substring(string, position);
Substring(string FROM position);
```

*   `string` 参数是要提取子字符串的字符串
*   `position` 参数是**正整数**，用于指定子串的起始字符， 从`1`开始，可以负整数

**使用方法：**

    SELECT SUBSTRING('MYSQL SUBSTRING', 7);

    结果是
    SUBSTRING

## 还可以指定长度

    SUBSTRING(string, position, length);
    SUBSTRING(string FROM position FOR length)

## 类似还有 `CONCAT、SUBSTR、SUBSTRING、SUBSTRING_INDEX、LEFT、RIGHT`

## 应用实例

    raw_sql = '''
                SELECT  SUBSTRING_INDEX(area_54350,",", :province_index) AS province, COUNT(object_id_48066) AS sale
                FROM tb_device 
                GROUP BY SUBSTRING_INDEX(area_54350,",", :province_index)
                HAVING province != ""
                ORDER BY sale DESC 
            '''
        params = {
            "province_index": 1,
        }
        fetchall = helper.get_prepared_query(raw_sql, params=params)
        provinces_data = []
        for row in fetchall:
            device = {
                "province": row[0],
                "count": row[1]
            }
            provinces_data.append(device)

    def get_prepared_query(raw_sql, params):
        stmt = text(raw_sql)
        return db.session.execute(stmt, params).fetchall()

