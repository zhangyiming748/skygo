DROP sequence IF EXISTS "public"."data_output_kafka_id_seq";

CREATE SEQUENCE "public"."data_output_kafka_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;


-- ----------------------------
-- Table structure for data_output_kafka
-- ----------------------------
DROP TABLE IF EXISTS "public"."data_output_kafka";
CREATE TABLE "public"."data_output_kafka" (
  "id" int4 NOT NULL DEFAULT nextval('data_output_kafka_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "config_kafka_id" int4 NOT NULL,
  "topics" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "output_data_format" int2 NOT NULL DEFAULT 1,
  "basic_template_id" int4,
  "create_time" int4 NOT NULL,
  "uuid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "category" int2
)
;
COMMENT ON COLUMN "public"."data_output_kafka"."id" IS '自增ID';
COMMENT ON COLUMN "public"."data_output_kafka"."name" IS '数据源名称';
COMMENT ON COLUMN "public"."data_output_kafka"."config_kafka_id" IS '表cluster_kafka的记录主键ID';
COMMENT ON COLUMN "public"."data_output_kafka"."topics" IS 'kafka主题，多个用逗号分割';
COMMENT ON COLUMN "public"."data_output_kafka"."output_data_format" IS '数据输出格式，1json，默认1';
COMMENT ON COLUMN "public"."data_output_kafka"."basic_template_id" IS '数据模板';
COMMENT ON COLUMN "public"."data_output_kafka"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."data_output_kafka"."uuid" IS 'uuid';
COMMENT ON COLUMN "public"."data_output_kafka"."category" IS '分类：1预定义，2自定义';
COMMENT ON TABLE "public"."data_output_kafka" IS 'kafka数据源配置';

-- ----------------------------
-- Records of data_output_kafka
-- ----------------------------
INSERT INTO "public"."data_output_kafka" VALUES (1, '威胁数据', 1, 'Threat', 1, 42, 1614915021, '94d3419f-899b-4528-b96c-41d36650cf4b', 1);
INSERT INTO "public"."data_output_kafka" VALUES (2, '协议事件—SSL加密协商', 1, 'Ssl', 1, 41, 1614915021, '0ac043d6-5253-4b98-93da-c0c80fc14bcd', 1);
INSERT INTO "public"."data_output_kafka" VALUES (3, '协议事件—WEB访问', 1, 'Web', 1, 43, 1614915021, '77b35784-b038-4752-b3bb-2941269661cc', 1);
INSERT INTO "public"."data_output_kafka" VALUES (4, '协议事件—域名解析', 1, 'Dns', 1, 33, 1614915021, '50a6806b-44c7-4404-b555-8644cfe7bea7', 1);
INSERT INTO "public"."data_output_kafka" VALUES (5, '协议事件—工业关键事件', 1, 'Ics_Event', 1, 37, 1614915021, '614bd6d0-54a6-4642-90c6-0efe86ae421a', 1);
INSERT INTO "public"."data_output_kafka" VALUES (6, '协议事件—邮件行为', 1, 'Mail', 1, 39, 1614915021, 'b1c8fcb0-681f-4d08-8973-3630dc5d1bff', 1);
INSERT INTO "public"."data_output_kafka" VALUES (7, '流量—会话关闭', 1, 'Flow_End', 1, 36, 1614915021, 'e47eefc4-e5b4-4586-bb25-4e207d53ef12', 1);
INSERT INTO "public"."data_output_kafka" VALUES (8, '白名单数据', 1, 'White', 1, 44, 1614915021, 'dd0d112b-e680-468c-b572-3467d9f45abe', 1);
INSERT INTO "public"."data_output_kafka" VALUES (9, '资产', 1, 'Asset', 1, NULL, 1614915021, '6e38f00a-28bd-41cb-a827-50f1c86e6c88', 1);

-- ----------------------------
-- Primary Key structure for table data_output_kafka
-- ----------------------------
ALTER TABLE "public"."data_output_kafka" ADD CONSTRAINT "_copy_2_copy_2" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table data_output_kafka
-- ----------------------------
ALTER TABLE "public"."data_output_kafka" ADD CONSTRAINT "fk_data_output_kafka_basic_template_1" FOREIGN KEY ("basic_template_id") REFERENCES "public"."basic_template" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."data_output_kafka" ADD CONSTRAINT "fk_data_output_kafka_config_kafka_1" FOREIGN KEY ("config_kafka_id") REFERENCES "public"."config_kafka" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
