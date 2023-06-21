DROP sequence IF EXISTS "public"."data_input_kafka_id_seq";

CREATE SEQUENCE "public"."data_input_kafka_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

SELECT setval('"public"."data_input_kafka_id_seq"', 11, true);

-- ----------------------------
-- Table structure for data_input_kafka
-- ----------------------------
DROP TABLE IF EXISTS "public"."data_input_kafka";
CREATE TABLE "public"."data_input_kafka" (
  "id" int4 NOT NULL DEFAULT nextval('data_input_kafka_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "config_kafka_id" int4 NOT NULL,
  "topics" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "consumer_group" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "decode" int2 NOT NULL DEFAULT 0,
  "avro_type" int2,
  "avro_content" varchar(255) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "create_time" int4 NOT NULL,
  "uuid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."data_input_kafka"."id" IS '自增ID';
COMMENT ON COLUMN "public"."data_input_kafka"."name" IS '数据源名称';
COMMENT ON COLUMN "public"."data_input_kafka"."config_kafka_id" IS '表cluster_kafka的记录主键ID';
COMMENT ON COLUMN "public"."data_input_kafka"."topics" IS 'kafka主题，多个用逗号分割';
COMMENT ON COLUMN "public"."data_input_kafka"."consumer_group" IS '消费者组名称';
COMMENT ON COLUMN "public"."data_input_kafka"."decode" IS '解码器，0无， 1avro';
COMMENT ON COLUMN "public"."data_input_kafka"."avro_type" IS '解码器类型，1 http 2 file';
COMMENT ON COLUMN "public"."data_input_kafka"."avro_content" IS '解码器avro的http地址或者file文件路径';
COMMENT ON COLUMN "public"."data_input_kafka"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."data_input_kafka" IS '数据输入kafka配置';

-- ----------------------------
-- Records of data_input_kafka
-- ----------------------------
INSERT INTO "public"."data_input_kafka" VALUES (11, '资产上报', 1, 'Asset-in', 'iisop', 0, NULL, NULL, 1614915021, 'f854cdff-2abc-419d-a9ea-50deac72e1b3');
INSERT INTO "public"."data_input_kafka" VALUES (9, '流量-会话关闭上报', 1, 'Flow_End-in', 'iisop', 0, NULL, '', 1614915021, '3b79ec48-d83c-4249-b264-e6ebb5b6c8b1');
INSERT INTO "public"."data_input_kafka" VALUES (1, '中能融合-行为日志上报', 1, 'Zhongneng_behavior_category-in', 'iisop', 0, NULL, NULL, 1614915021, 'cc9f6216-a852-48c1-8336-0cc2b63e44a1');
INSERT INTO "public"."data_input_kafka" VALUES (2, '中能融合-警报日志上报', 1, 'Zhongneng_alarm_category-in', 'iisop', 0, NULL, NULL, 1614915021, '7e6274db-bce4-47c6-ad30-42459f25a87f');
INSERT INTO "public"."data_input_kafka" VALUES (3, '协议-SSL加密协商上报', 1, 'Ssl-in', 'iisop', 0, NULL, NULL, 1614915021, 'a925e5c2-f6a3-4a67-a821-6bd2329049c0');
INSERT INTO "public"."data_input_kafka" VALUES (4, 'WEB访问上报', 1, 'Web-in', 'iisop', 0, NULL, '', 1614915021, 'a8408ad5-96bf-4fa3-a6d7-8e257cbba001');
INSERT INTO "public"."data_input_kafka" VALUES (5, '协议-域名解析上报', 1, 'Dns-in', 'iisop', 0, NULL, '', 1614915021, 'd82174a9-5e0d-43a8-a145-78b3596ea3cf');
INSERT INTO "public"."data_input_kafka" VALUES (6, '协议-威胁数据上报', 1, 'Threat-in', 'iisop', 0, NULL, '', 1614915021, '563de29a-ebde-4cc0-a1ab-f5af80d89993');
INSERT INTO "public"."data_input_kafka" VALUES (7, '协议-工业关键事件上报', 1, 'Ics_Event-in', 'iisop', 0, NULL, '', 1614915021, '288b76d5-fbc6-4aa9-99ba-ad8e8eed200b');
INSERT INTO "public"."data_input_kafka" VALUES (8, '协议-邮件行为上报', 1, 'Mail-in', 'iisop', 0, NULL, '', 1614915021, '856cb227-de32-4509-8a55-c0bbd95b65d5');
INSERT INTO "public"."data_input_kafka" VALUES (10, '白名单数据上报', 1, 'Mail-in', 'iisop', 0, NULL, '', 1614915021, 'fcf97c35-d229-446a-841f-5561dee59aab');

-- ----------------------------
-- Uniques structure for table data_input_kafka
-- ----------------------------
ALTER TABLE "public"."data_input_kafka" ADD CONSTRAINT "data_input_kafka.uuid" UNIQUE ("uuid");

-- ----------------------------
-- Primary Key structure for table data_input_kafka
-- ----------------------------
ALTER TABLE "public"."data_input_kafka" ADD CONSTRAINT "_copy_2" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table data_input_kafka
-- ----------------------------
ALTER TABLE "public"."data_input_kafka" ADD CONSTRAINT "fk_data_input_kafka_config_kafka_1" FOREIGN KEY ("config_kafka_id") REFERENCES "public"."config_kafka" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
