# Usage of weibo2emo:
#### -db 
*是否生成mysql数据库，默认为 false (选填)*

#### -gf string 
*分级时所使用的映射函数，目前可选的有:*
        
> linear : f(x) = x
> 
> log : f(x) = ln(x + 1)
>
> sigmoid : f(x) = 1 / (1 + exp(-x)) + 0.5
>
> tanh : 2 / (1 + exp(-2x)) -1

*默认使用线性函数 (选填) (default "linear")*

#### -gn int
*所要分出的等级数量，默认为 5 (选填) (default 5)*

#### -gp string
*待分级的数据集(.csv)所在路径 示例 ./result_time.csv (如要进行分级统计，则必填，否则，非必填)*
#### -gr string
*获取到的分级结果，默认导出在当前目录，名称为result_grade.csv，如想自定义导出，一定要精确到文件名 (选填) (default "./result_grade.csv")*
#### -pd string
*字典集数据(.csv)所在的路径 示例 ./LIWC.csv (如要进行词频统计，则必填，否则，非必填)*
#### -pp string
*博文数据(.csv)所在的路径 示例 ./原始博文1月.csv (如要进行词频统计，则必填，否则，非必填)*
#### -pr string
*结果数据（未按时间聚合）想要导出的目录，默认导出在当前目录，名称为result.csv,如想自定义导出，一定要精确到文件名 (选填) (default "./result.csv")*
#### -prt string
*按时间聚合的结果数据集想要导出的目录，默认导出在当前目录，名称为result_time.csv，如想自定义导出，一定要精确到文 件名 (选填) (default "./result_time.csv")*
#### -pt string
*id,用户名数据集(.csv)所在路径 示例 ./博主名称及id.csv (如要进行词频统计，则必填，否则，非必填)*
#### -rh int
*待输出的图片高度，默认为40，非必填 (default 40)*
#### -rp string
*待输出绘制的csv文件所在位置 实例./result_time.csv (如要进行折线绘制则必填，否则，非必填)*
#### -rw int
*待输出的图片长度，默认为150，非必填 (default 150)*
#### -thread int
*使用的cpu线程数 (default 5)*
