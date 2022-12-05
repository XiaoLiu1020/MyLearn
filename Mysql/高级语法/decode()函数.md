# decode

作用：获取值，符合条件，则返回改变值

    decode (item, if1, then1, if2, then2,...,else)

## 可以使用`decode`比较大小

    select decode(sign(var1-va2), -1, var1, var2) from dual
    // sign()函数根据某个值是0，正数还是负数，分别返回0，1，-1

    select decode(sign(100-90)), -1, 100, 90) from dual
    90
    //100-90=10>0，则返回1，所以decode函数最终取值为90

    select decode(sign(100-90)), 1, 100, 90) from dual
    100
    100-90=10>0返回1，结果为1，返回第一个变量，最终为100

## 使用`decode`函数分段

## `mysql`并没有`decode`函数，只能换种写法,使用`case when`

```
select decode(pay_name,'aaaa','bbb',pay_name),sum(comm_order),sum(suc_order),sum(suc_amount) From  payment.order_tab  group by decode(pay_name,'aaaaa','bbbb',pay_name)

// 转换成mysql:实现 ,使用 case when
 
select case when pay_name='aaa' then 'bbb' else pay_name end ,sum(comm_order),sum(suc_order),sum(suc_amount) From  payment.order_tab  group by case when pay_name='aaa' then 'bbb' else pay_name end 

```

