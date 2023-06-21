DROP sequence IF EXISTS "public"."basic_template_id_seq";

CREATE SEQUENCE "public"."basic_template_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

SELECT setval('"public"."basic_template_id_seq"', 44, true);


-- ----------------------------
-- Table structure for basic_template
-- ----------------------------
DROP TABLE IF EXISTS "public"."basic_template";
CREATE TABLE "public"."basic_template" (
  "id" int4 NOT NULL DEFAULT nextval('basic_template_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "alias" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "if_kafka_memory" int2 NOT NULL DEFAULT 0,
  "if_es_memory" int2 NOT NULL DEFAULT 0,
  "es_searchable" int2 NOT NULL DEFAULT 0,
  "es_partition_mode" int2 NOT NULL DEFAULT 0,
  "es_minimum_day" int4 NOT NULL DEFAULT 31,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "category" int2 NOT NULL DEFAULT 2,
  "create_time" int4 NOT NULL DEFAULT 0,
  "es_index_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "log_group_id" int4 DEFAULT 0
)
;
COMMENT ON COLUMN "public"."basic_template"."id" IS '自增主键';
COMMENT ON COLUMN "public"."basic_template"."name" IS '数据模板的名称';
COMMENT ON COLUMN "public"."basic_template"."alias" IS '数据模板的别名';
COMMENT ON COLUMN "public"."basic_template"."if_kafka_memory" IS '是否用于kafka存储模板，0否， 1是';
COMMENT ON COLUMN "public"."basic_template"."if_es_memory" IS '是否创建es索引， 0否， 1是';
COMMENT ON COLUMN "public"."basic_template"."es_searchable" IS '是否用于es检索，0否，1是';
COMMENT ON COLUMN "public"."basic_template"."es_partition_mode" IS 'es分区类型，0不分区， 1按天分区';
COMMENT ON COLUMN "public"."basic_template"."es_minimum_day" IS 'es如果按时间分区，最低保存天数，默认31
';
COMMENT ON COLUMN "public"."basic_template"."description" IS '描述';
COMMENT ON COLUMN "public"."basic_template"."category" IS '分类：1预定义，2自定义';
COMMENT ON COLUMN "public"."basic_template"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."basic_template"."es_index_name" IS '索引名称';
COMMENT ON COLUMN "public"."basic_template"."log_group_id" IS '日志组id，支持es索引的，可以设置日志组id';

-- ----------------------------
-- Records of basic_template
-- ----------------------------
INSERT INTO "public"."basic_template" VALUES (37, 'ics_event', '工业关键事件', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_ics_event', 6);
INSERT INTO "public"."basic_template" VALUES (33, 'dns', 'Dns请求', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_dns', 5);
INSERT INTO "public"."basic_template" VALUES (38, 'login', '登入日志', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_login', 5);
INSERT INTO "public"."basic_template" VALUES (39, 'mail', '邮件信息', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_mail', 5);
INSERT INTO "public"."basic_template" VALUES (40, 'sql', '数据库操作', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_sql', 5);
INSERT INTO "public"."basic_template" VALUES (41, 'ssl', 'SSL协议', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_ssl', 5);
INSERT INTO "public"."basic_template" VALUES (43, 'web', 'Web访问', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_web', 5);
INSERT INTO "public"."basic_template" VALUES (34, 'file', '文件上传事件', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_file', 4);
INSERT INTO "public"."basic_template" VALUES (44, 'white', '白名单事件', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_white', 3);
INSERT INTO "public"."basic_template" VALUES (42, 'threat', '威胁日志', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_threat', 1);
INSERT INTO "public"."basic_template" VALUES (35, 'flow_begin', '流量开始', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_flow_begin', 2);
INSERT INTO "public"."basic_template" VALUES (32, 'alarm_event', '威胁告警', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_alarm_event', 1);
INSERT INTO "public"."basic_template" VALUES (1, 'nanwang_general', '通用日志', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_nanwang_general', 7);
INSERT INTO "public"."basic_template" VALUES (3, 'zhongneng_behavior', '中能行为日志', 1, 1, 1, 1, 180, '', 2, 1614684337, 'iisop_zhongneng_behavior', 8);
INSERT INTO "public"."basic_template" VALUES (36, 'flow_end', '流量结束', 1, 1, 1, 1, 30, '', 2, 1614684337, 'iisop_flow_end', 2);

-- ----------------------------
-- Primary Key structure for table basic_template
-- ----------------------------
ALTER TABLE "public"."basic_template" ADD CONSTRAINT "_copy_1_copy_1_copy_1" PRIMARY KEY ("id");
