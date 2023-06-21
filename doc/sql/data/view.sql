-- 视图
-- 数据模板,为了方便查询引用次数,使用视图
CREATE VIEW basic_template_view AS
SELECT bt.*,
    (COALESCE(es.total, 0) + COALESCE(kafka.total,0)) AS quote_total
FROM ((basic_template bt
LEFT JOIN
(
	SELECT basic_template_id AS id, count(1) AS total
	FROM data_output_es
	WHERE basic_template_id != 0 AND basic_template_id IS NOT NULL
	GROUP BY basic_template_id
) es ON ((bt.id = es.id)))
LEFT JOIN
(
	SELECT basic_template_id AS id, count(1) AS total
	FROM data_output_kafka
	WHERE basic_template_id != 0 AND basic_template_id IS NOT NULL
	GROUP BY basic_template_id
) kafka ON ((bt.id = kafka.id)));

-- kafka基础存储配置,为了方便查询使用次数,使用视图
create view config_kafka_view as
select c.*, (count(dok) + count(dik)) as quote_total
from config_kafka as c
left join data_output_kafka as dok on dok.config_kafka_id = c.id
left join data_input_kafka as dik on dik.config_kafka_id = c.id
group by c.id;

-- es基础存储配置
create view config_es_view as
select c.*, count(doe) as quote_total
from config_es as c
left join data_output_es as doe on doe.config_es_id = c.id
group by c.id;

--- 数据输入
create view data_output_view as
select  'kafka' AS "output_type","uuid","name",create_time from data_output_kafka
union all
select 'es' AS "output_type","uuid","name",create_time from data_output_es;

-- 数据存储es 视图
create view data_output_es_view as
select e.*, count(o) as quote_total
from data_output_es as e
left join task_data_output as o
on o.data_output_uuid = e.uuid
group by e.id;

-- 数据存储kafka 视图
create view data_output_kafka_view as
select e.*, count(o) as quote_total
from data_output_kafka as e
left join task_data_output as o
on o.data_output_uuid = e.uuid
group by e.id;

-- 基础字段视图
-- 基础字段表连接基础模板字段表, 能知道一个字段被引用的次数
create view basic_field_view as
select bf.*, count(btf) as quote_total
from basic_field as bf
left join basic_template_field as btf on btf.field_id = bf."id"
group by bf.id;


-- 任务视图
-- 连接rule表得到解析规则名称
-- 连接myin(所有输入源的数据整合)得到输入源名称
-- 连接task_out（所有输出源的数据整合）得到输出的名称
CREATE VIEW task_view AS
SELECT t.*, task_out.*, myin."name" as input_name, r."name" as rule_name
FROM task as t
LEFT JOIN rule as r ON t.rule_id = r.id
LEFT JOIN
(
	SELECT uuid, name FROM data_input_kafka
	UNION
	SELECT uuid, name FROM data_input_syslog
) AS myin ON t.data_input_uuid = myin.uuid
LEFT JOIN
(
	SELECT
	tao.task_id,
	string_agg(tao.data_output_uuid, ',') as output_uuids,
	string_agg(uout."name", ',') as output_names
	FROM task_data_output as tao
	LEFT JOIN
	(
	SELECT uuid, name, 'kafka' as class  FROM data_output_kafka
	UNION
	SELECT uuid, name, 'es' as class FROM data_output_es
	) as uout
	ON tao.data_output_uuid = uout.uuid
	GROUP BY tao.task_id
) AS task_out ON task_out.task_id = t.id;

--- 数据输入kafka 视图
create view data_input_kafka_view as
select kafka.*,config_kafka."name" AS config_kafka_name from (select i.*, count(t) as quote_total
from data_input_kafka as i
left join task as t
on i."uuid" = t.data_input_uuid
group by i.id) AS kafka LEFT JOIN config_kafka ON kafka.config_kafka_id=config_kafka.id;

--- 数据输入syslog 视图
create view data_input_syslog_view as
select syslog.*, count(task) as quote_total
from data_input_syslog as syslog
left join task
on syslog."uuid" = task.data_input_uuid
group by syslog.id;

--- 数据输入 视图
create view data_input_view as
select  'kafka' AS "input_type","uuid","name",id,create_time from data_input_kafka
union all
select 'syslog' AS "input_type","uuid","name",id,create_time from data_input_syslog;

--- 规则 视图
CREATE VIEW rule_view AS
SELECT r.*, count(t) as quote_total
FROM rule AS r
LEFT JOIN task AS t ON t.rule_id = r."id"
GROUP BY r."id";

--- 模板字段视图
--- 方便查询模板字段时，也能拿到字段名称、字段类型等信息，从而不用在连表查询
CREATE VIEW basic_template_field_view AS
SELECT v.*, f."name" AS "name", f."type" as "type", f."alias" as "alias"
FROM basic_template_field AS v
LEFT JOIN basic_field AS f ON v.field_id = f."id";