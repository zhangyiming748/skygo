DROP sequence IF EXISTS "public"."data_input_syslog_id_seq";

CREATE SEQUENCE "public"."data_input_syslog_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;


-- ----------------------------
-- Table structure for data_input_syslog
-- ----------------------------
DROP TABLE IF EXISTS "public"."data_input_syslog";
CREATE TABLE "public"."data_input_syslog" (
  "id" int4 NOT NULL DEFAULT nextval('data_input_syslog_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "service_port" int4 NOT NULL DEFAULT 0,
  "character_encoding" int2 NOT NULL DEFAULT 1,
  "source_ips" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "create_time" int4 NOT NULL,
  "uuid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
COMMENT ON COLUMN "public"."data_input_syslog"."id" IS '自增ID';
COMMENT ON COLUMN "public"."data_input_syslog"."name" IS '数据源名称';
COMMENT ON COLUMN "public"."data_input_syslog"."service_port" IS '绑定端口';
COMMENT ON COLUMN "public"."data_input_syslog"."character_encoding" IS '字符编码， 1utf-8';
COMMENT ON COLUMN "public"."data_input_syslog"."source_ips" IS '报送设备的ip地址，格式ip,ip,ip';
COMMENT ON COLUMN "public"."data_input_syslog"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."data_input_syslog" IS '数据输入syslog配置';

-- ----------------------------
-- Uniques structure for table data_input_syslog
-- ----------------------------
ALTER TABLE "public"."data_input_syslog" ADD CONSTRAINT "data_input_syslog.uuid" UNIQUE ("uuid");

-- ----------------------------
-- Primary Key structure for table data_input_syslog
-- ----------------------------
ALTER TABLE "public"."data_input_syslog" ADD CONSTRAINT "_copy_2_copy_1" PRIMARY KEY ("id");
